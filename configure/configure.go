package configure

import (
	"strconv"
	"time"

	"github.com/nece099/base/encrypt"
)

type Model struct {
	ID        int64 `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Configure struct {
	Model
	ParamName  string `gorm:"size:64"`
	ParamValue string `gorm:"type:longtext"`
	Disabled   bool   `gorm:"default:0"`
	Encrypted  bool   `gorm:"default:0"`
	Remark     string `gorm:"type:longtext"`
}

func (c *Configure) String() string {
	return c.ParamValue
}

func (c *Configure) Int64() int64 {

	sval := c.ParamValue
	i, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func (c *Configure) Float64() float64 {

	sval := c.ParamValue
	f, err := strconv.ParseFloat(sval, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (c *Configure) Decrypt() string {

	d, err := encrypt.InternalDecryptStr(c.String())
	if err != nil {
		Log.Errorf("decrypt failed, err = %v", err)
		return ""
	}

	return d
}

func (c *Configure) IsDisabled() bool {
	return c.Disabled
}
