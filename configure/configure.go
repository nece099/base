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

type Item struct {
	Value     interface{}
	Disabled  bool
	Encrypted bool
}

func (item *Item) String() string {
	return item.Value.(string)
}

func (item *Item) Int64() int64 {

	sval := item.Value.(string)
	i, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func (item *Item) Float64() float64 {

	sval := item.Value.(string)

	f, err := strconv.ParseFloat(sval, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func (item *Item) Decrypt() string {

	d, err := encrypt.InternalDecryptStr(item.String())
	if err != nil {
		Log.Errorf("decrypt failed, err = %v", err)
		return ""
	}

	return d
}

func (item *Item) IsDisabled() bool {
	return item.Disabled
}
