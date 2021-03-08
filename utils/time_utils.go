package utils

import (
	"fmt"
	"time"
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

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

/////////////////////////////////////////////////////////////
func NowFormatYMDHMS() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func DateStr(t time.Time) string {
	ds := fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
	return ds
}

func DateParse(dstr string) (time.Time, error) {
	return time.Parse(DATE_FORMAT, dstr)
}
