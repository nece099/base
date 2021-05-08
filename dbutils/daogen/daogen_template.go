package daogen

var daogen_template = `
package dao

import (
	"errors"
	"sync"

	"gorm.io/gorm"
	"github.com/nece099/base/except"
	"github.com/nece099/base/dbo"
	"%v/do"
	"github.com/nece099/base/dbutils"
)

type {{.StructName}}Dao struct {
	mutex *sync.Mutex
}

var {{LowerCaseFirstLetter .StructName}}Dao *{{.StructName}}Dao = nil

func New{{.StructName}}Dao() *{{.StructName}}Dao {
	{{LowerCaseFirstLetter .StructName}}Dao = &{{.StructName}}Dao{
		mutex: &sync.Mutex{},
	}
	return {{LowerCaseFirstLetter .StructName}}Dao
}

func Get{{.StructName}}Dao() *{{.StructName}}Dao {
	except.ASSERT({{LowerCaseFirstLetter .StructName}}Dao != nil)
	return {{LowerCaseFirstLetter .StructName}}Dao
}

func (dao *{{.StructName}}Dao) Lock() {
	dao.mutex.Lock()
}

func (dao *{{.StructName}}Dao) Unlock() {
	dao.mutex.Unlock()
}

func (dao *{{.StructName}}Dao)DB() *gorm.DB {
	return dbo.DboInstance().DB()
}


func (dao *{{.StructName}}Dao) Create(m *do.{{.StructName}}) error {
	return dao.DB().Create(m).Error
}

func (dao *{{.StructName}}Dao) Find(m *do.{{.StructName}}) (result []*do.{{.StructName}}, err error) {
	err = dao.DB().Find(&result, m).Error
	return
}

func (dao *{{.StructName}}Dao) FindOne(m *do.{{.StructName}}) error {
	return dao.DB().First(m, m).Error
}

func (dao *{{.StructName}}Dao) FindLast(m *do.{{.StructName}}) error {
	return dao.DB().Last(m, m).Error
}

func (dao *{{.StructName}}Dao) FindPage(m *do.{{.StructName}}, p *dbutils.Paging) (result []*do.{{.StructName}}, err error) {

	db := dao.DB()

	db = db.Model(m).Where(m)

	if p == nil {
		err = db.Find(&result).Error
		if err != nil {
			return
		}
	} else {

		err = db.Count(&p.AllCount).Error
		if err != nil {
			return
		}

		err = db.Limit(int(p.PageSize)).
			Offset(int(p.PageIndex * p.PageSize)).
			Find(&result).Error
		if err != nil {
			return
		}
	}

	return
}

func (dao *{{.StructName}}Dao) BatchGet(idbatch []int64) (result []*do.{{.StructName}}, err error) {
	if len(idbatch) == 0 {
		return nil, errors.New("id is nil")
	}
	err = dao.DB().Model(&do.{{.StructName}}{}).Where("id in (?)", idbatch).Find(&result).Error
	return
}

func (dao *{{.StructName}}Dao) Save(m *do.{{.StructName}}) error {
	return dao.DB().Save(m).Error
}

func (dao *{{.StructName}}Dao) BatchDelete(idbatch []int64) error {
	if len(idbatch) == 0 {
		return errors.New("id is nil")
	}
	return dao.DB().Where("id in (?)", idbatch).Delete(&do.{{.StructName}}{}).Error
}

func (dao *{{.StructName}}Dao) Updates(id int64, attrs map[string]interface{}) error {
	return dao.DB().Model(&do.{{.StructName}}{}).Where("id = ?", id).Updates(attrs).Error
}

func (dao *{{.StructName}}Dao) Update(id int64, attr string, value interface{}) error {
	return dao.DB().Model(&do.{{.StructName}}{}).Where("id = ?", id).Update(attr, value).Error
}

func (dao *{{.StructName}}Dao) BatchUpdaterAttrs(idbatch []int64, attrs map[string]interface{}) error {
	if len(idbatch) == 0 {
		return errors.New("id is nil")
	}
	return dao.DB().Model(&do.{{.StructName}}{}).Where("id in (?)", idbatch).Updates(attrs).Error
}

func (dao *{{.StructName}}Dao) Found(m *do.{{.StructName}}) bool {
	err := dao.db.First(m, m).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

`
