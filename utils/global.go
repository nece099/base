package utils

import (
	"bytes"
	"encoding/hex"
	json2 "encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/leekchan/accounting"
	. "github.com/nece099/base/logger"
)

const (
	TIME_FORMAT_WITH_MS         = "2006-01-02 15:04:05.000"
	TIME_FORMAT                 = "2006-01-02 15:04:05"
	TIME_FORMAT_COMPACT         = "20060102150405"
	TIME_FORMAT_WITH_MS_COMPACT = "20060102150405.000"
	DATE_FORMAT                 = "2006-01-02"
	DATE_FORMAT_COMPACT         = "20060102"
	MONTH_FORMAT                = "2006-01"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NowDateStr() string {
	timenow := time.Now().Format(DATE_FORMAT)
	return timenow
}

func UseMaxCpu() {
	// multiple cups using
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func FormatJson(instr string) string {

	var out bytes.Buffer
	json2.Indent(&out, []byte(instr), "", "  ")

	return "\n" + out.String() + "\n"
}

func FormatStruct(inst interface{}) string {
	instr := SerializeToJson(inst)
	return FormatJson(instr)
}

func SerializeToJson(st interface{}) string {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(st)

	return buf.String()
}

// instead of original unserialize function
func UnserializeFromJson(jsonstr string, st interface{}) error {
	d := json.NewDecoder(strings.NewReader(jsonstr))
	d.UseNumber()
	return d.Decode(st)
}

func UnserializeJson(jsonstr string, st interface{}) {
	d := json.NewDecoder(strings.NewReader(jsonstr))
	d.UseNumber()
	err := d.Decode(st)
	if err != nil {
		Log.Error("parse json failed, err = %v", err)
		panic(err)
	}
}

func HexBuffer(buffer []byte) string {
	s := hex.EncodeToString(buffer)

	n := 8
	m := 8

	c := 0
	slen := len(s)
	if slen%n == 0 {
		c = slen / n
	} else {
		c = (slen / n) + 1
	}

	res := ""
	for i := 0; i < c; i++ {
		res = res + s[i:i+8] + " "
		if (i+1)%m == 0 {
			res = res + "\n"
		}
	}

	res = fmt.Sprintf("\n=======================================================================\n%v\n=======================================================================\n", res)
	return res
}

func GetProgName() string {
	fullPath, _ := exec.LookPath(os.Args[0])
	fname := filepath.Base(fullPath)

	return fname
}

func EncodeURI(data string) string {
	return url.QueryEscape(data)
}

func DecodeURI(data string) (string, error) {

	sdata, err := url.QueryUnescape(data)
	if err != nil {
		Log.Errorf("url.QueryUnescape err = %v", err)
		return "", err
	}

	return sdata, nil
}

func IsNullTime(t time.Time) bool {
	year := t.Year()
	if year == 1 {
		return true
	} else {
		return false
	}
}

// IsLower check letter is lower case or not
func IsLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

// IsUpper check letter is upper case or not
func IsUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// IsLetter check character is a letter or not
func IsLetter(b byte) bool {
	return IsLower(b) || IsUpper(b)
}

// IsNumber check character is a number or not
func IsNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsAlNum check character is a alnum or not
func IsAlNum(c byte) bool {
	return IsLetter(c) || IsLetter(c)
}

func TimeStr(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func FormatMoneyStr(in string) string {
	infloat64 := String2Float64(in)
	ac := accounting.Accounting{Symbol: "", Precision: 3}
	fmoney := ac.FormatMoney(infloat64)
	return fmoney
}
