package configure

import (
	"strconv"
	"sync"

	"time"

	"github.com/nece099/base/dbo"
	"github.com/nece099/base/encrypt"
	"github.com/nece099/base/except"
	"github.com/nece099/base/grmon"
)

type Model struct {
	ID        int64 `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProgConfig struct {
	Model
	ParamName  string `gorm:"size:64"`
	ParamValue string `gorm:"type:longtext"`
	Type       string `gorm:"size:32"`
	Disabled   bool   `gorm:"default:0"`
	Encrypted  bool   `gorm:""`
	Comment    string `gorm:"size:1024"`
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

type ProgConfigure struct {
	configMap *sync.Map
}

var progConfigure *ProgConfigure = &ProgConfigure{
	configMap: &sync.Map{},
}

// var progConfigMap cmap.ConcurrentMap = cmap.New()

func ProgConfigureInit(defaultConfig map[string]string) {
	setConfig(defaultConfig)
	loadProgConfig()

	grm := grmon.GetGRMon()
	grm.Go("reloadProgConfig", func() {
		for {
			time.Sleep(time.Duration(60) * time.Second)
			loadProgConfig()
		}
	})
}

func setConfig(configmap map[string]string) {

	for k, v := range configmap {

		item := &Item{
			Value:    v,
			Disabled: false,
		}

		progConfigure.configMap.Store(k, item)
	}
}

func loadProgConfig() {
	var configs []*ProgConfig

	db := dbo.DboInstance().DB()
	err := db.Raw("/*no print*/ select * from prog_config").Find(&configs).Error
	if err != nil {
		Log.Warnf("load prog config failed, err = %v", err)
		return
	}

	for _, c := range configs {

		item := &Item{
			Value:     c.ParamValue,
			Disabled:  c.Disabled,
			Encrypted: c.Encrypted,
		}
		progConfigure.configMap.Store(c.ParamName, item)
	}
}

func ProgConfigureInstance() *ProgConfigure {
	except.ASSERT(progConfigure != nil)
	return progConfigure
}

func (p *ProgConfigure) GetConfigItem(name string) *Item {
	itemobj, ok := p.configMap.Load(name)
	except.ASSERT(ok)

	return itemobj.(*Item)
}

func (p *ProgConfigure) GetConfigItemFromDb(name string) *Item {
	config := ProgConfig{}

	db := dbo.DboInstance().DB()
	err := db.Where("param_name=?", name).Find(&config).Error
	if err != nil {
		Log.Error("db error = %v", err)
		panic(err)
	}

	return &Item{Value: config.ParamValue, Disabled: config.Disabled}
}

func (p *ProgConfigure) SetConfigItemWithEncryption(name string, value string) error {

	encrypted, err := encrypt.InternalEncryptStr(value)
	if err != nil {
		Log.Errorf("value=%v, err=%v", value, err)
		return err
	}

	return p.SetConfigItem(name, encrypted)
}

func (p *ProgConfigure) SetConfigItem(name string, value string) error {

	config := &ProgConfig{
		ParamName:  name,
		ParamValue: value,
	}

	db := dbo.DboInstance().DB()
	err := db.Save(config).Error
	if err != nil {
		Log.Errorf("db error = %v", err)
		return err
	}

	p.configMap.Store(name, config)

	return nil
}

func GetConfigItem(name string) *Item {
	return ProgConfigureInstance().GetConfigItem(name)
}

func SetConfigItem(name string, value string) error {
	return ProgConfigureInstance().SetConfigItem(name, value)
}

func GetConfigItemDb(name string) *Item {
	return ProgConfigureInstance().GetConfigItemFromDb(name)
}

func SetConfigItemWithEncryption(name string, value string) error {
	return ProgConfigureInstance().SetConfigItemWithEncryption(name, value)
}
