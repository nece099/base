package daogen

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nece099/base/logger"
	"github.com/nece099/base/utils"
)

var Log *logger.Logger = nil

func init() {
	Log = logger.Log
}

var dao_template *template.Template

type ModelFill struct {
	StructLines string
	InitLines   string
}

func DaoGenEntry() {
	input := ""
	output := ""
	modelPath := ""
	modelPackge := ""

	flag.StringVar(&input, "i", "./model/do", "input files")
	flag.StringVar(&modelPath, "m", "./model/model.go", "model.go")
	flag.StringVar(&output, "o", "./model/dao", "output directory")
	flag.StringVar(&modelPackge, "p", "", "model package")
	flag.Parse()

	if len(modelPackge) == 0 {
		Log.Warnf("model package must be specified")
		os.Exit(-1)
	}

	var fpaths []string

	files, _ := ioutil.ReadDir(input)
	for _, fi := range files {
		if !fi.IsDir() {
			path := input + "/" + fi.Name()
			fpaths = append(fpaths, path)
		}
	}

	// Log.Debugf("fpaths = %v", fpaths)
	var err error
	templateContent := fmt.Sprintf(daogen_template, modelPackge)
	dao_template, err = template.New("dao").Funcs(template.FuncMap{
		"LowerCaseFirstLetter": utils.LowerCaseFirstLetter,
	}).Parse(string(templateContent))
	if err != nil {
		Log.Error(err)
		return
	}

	modelLine1 := ""
	modelLine2 := ""

	for _, fpath := range fpaths {
		base := filepath.Base(fpath)
		if base != "do.go" &&
			base != "init.go" {
			sti := do2dao(fpath, output)
			if len(sti.StructName) == 0 {
				continue
			}

			daoName := fmt.Sprintf("%vDao *dao.%vDao\n", sti.StructName, sti.StructName)
			daoNew := fmt.Sprintf("model.%vDao = dao.New%vDao()\n", sti.StructName, sti.StructName)
			modelLine1 = modelLine1 + daoName
			modelLine2 = modelLine2 + daoNew
		}
	}

	//
	mf := &ModelFill{
		StructLines: modelLine1,
		InitLines:   modelLine2,
	}
	writeModel(modelPath, modelPackge, mf)

	// write dao
	writDao(output)
}

func writeModel(modelPath string, modelPackge string, mf *ModelFill) {
	modelgo := fmt.Sprintf(model_template, modelPackge, mf.StructLines, mf.InitLines)
	ioutil.WriteFile(modelPath, []byte(modelgo), 0666)
}

func writDao(output string) {
	ioutil.WriteFile(output, []byte(daogo_template), 0666)
}

type StructInfo struct {
	StructName string
	FieldNames []string
}

func do2dao(dofile string, outpath string) *StructInfo {

	fileContent, err := ioutil.ReadFile(dofile)
	if err != nil {
		Log.Error(err)
		return nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(fileContent), parser.AllErrors)

	// spew.Dump(f)

	sti := &StructInfo{}

	ast.Inspect(f, func(n ast.Node) bool {
		// Log.Debugf("n type = %v", reflect.TypeOf(n))
		switch n.(type) {

		case *ast.TypeSpec:
			r := n.(*ast.TypeSpec)
			// Log.Debugf("r spec = %v", r.Name)
			sti.StructName = r.Name.Name
		case *ast.StructType:
			r := n.(*ast.StructType)
			// Log.Debugf("r = %#v, pos = %v, end = %v", r, r.Pos(), r.End())
			// Log.Debugf("r.Fields = %#v", r.Fields)
			for _, v := range r.Fields.List {
				nlen := len(v.Names)
				if nlen > 0 {
					// Log.Debugf("v = %#v", v.Names[0].Name)
					sti.FieldNames = append(sti.FieldNames, v.Names[0].Name)
				}
			}
		}

		return true
	})

	// 生成文件
	outFile := outpath + "/" + utils.ToSnakeCase(sti.StructName) + "_dao.go"
	Log.Debugf("outFile = %v", outFile)
	ofile, err := os.OpenFile(outFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		Log.Error(err)
		return nil
	}

	err = dao_template.Execute(ofile, sti)
	if err != nil {
		Log.Error(err)
		return nil
	}

	ofile.Close()

	return sti
}
