package utils

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func LowerCaseFirstLetter(str string) string {
	if len(str) > 0 {
		first := strings.ToLower(string(str[0]))
		return first + string(str[1:])
	}
	return str
}

func StringToInt32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

func StringToUint32(s string) (uint32, error) {
	i, err := strconv.Atoi(s)
	return uint32(i), err
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func StringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func BoolToInt8(b bool) int8 {
	if b {
		return 1
	} else {
		return 0
	}
}

func Int8ToBool(b int8) bool {
	if b == 1 {
		return true
	} else {
		return false
	}
}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func JoinInt32List(s []int32, sep string) string {
	l := len(s)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(sep)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = strconv.AppendInt(buf, int64(s[i]), 10)
		// buf = append(buf, sep...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func JoinUint32List(s []uint32, sep string) string {
	l := len(s)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(sep)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = strconv.AppendUint(buf, uint64(s[i]), 10)
		buf = append(buf, sep...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func JoinInt64List(s []int64, sep string) string {
	l := len(s)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(sep)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = strconv.AppendInt(buf, s[i], 10)
		buf = append(buf, sep...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func JoinUint64List(s []uint64, sep string) string {
	l := len(s)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(sep)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = strconv.AppendUint(buf, s[i], 10)
		buf = append(buf, sep...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func String2Float64(num string) float64 {
	f, err := strconv.ParseFloat(num, 64)
	ASSERT(err == nil)

	return f
}

// IsAlNumString returns true if an alpha numeric string consists of characters a-zA-Z0-9
func IsAlNumString(s string) bool {
	c := 0
	for _, r := range s {
		switch {
		case '0' <= r && r <= '9':
			c++
			break
		case 'a' <= r && r <= 'z':
			c++
			break
		case 'A' <= r && r <= 'Z':
			c++
			break
		}
	}
	return len(s) == c
}

func RemoveDuplicate(arr []string) []string {
	result := make([]string, 0, len(arr))
	temp := map[string]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
