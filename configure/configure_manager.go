package configure

import (
	"sync"

	"time"

	"github.com/nece099/base/dbo"
	"github.com/nece099/base/encrypt"
	"github.com/nece099/base/except"
	"github.com/nece099/base/grmon"
)

type ConfigureManager struct {
	cache *sync.Map
}

var configureManager *ConfigureManager = &ConfigureManager{
	cache: &sync.Map{},
}

// var progConfigMap cmap.ConcurrentMap = cmap.New()

func ConfigureManagerInit(externalConfigure map[string]string) {
	setConfigures(externalConfigure)
	loadConfigures()

	grm := grmon.GetGRMon()
	grm.Go("reloadConfigures", func() {
		for {
			time.Sleep(time.Duration(60) * time.Second)
			loadConfigures()
		}
	})
}

func setConfigures(configMap map[string]string) {
	for k, v := range configMap {
		item := &Item{
			Value:    v,
			Disabled: false,
		}

		configureManager.cache.Store(k, item)
	}
}

func loadConfigures() {
	var cs []*Configure

	db := dbo.DboInstance().DB()
	err := db.Raw("/*no print*/ select * from configure").Find(&cs).Error
	if err != nil {
		Log.Warnf("load configure failed, err = %v", err)
		return
	}

	for _, c := range cs {

		item := &Item{
			Value:     c.ParamValue,
			Disabled:  c.Disabled,
			Encrypted: c.Encrypted,
		}
		configureManager.cache.Store(c.ParamName, item)
	}
}

func ConfigureManagerInstance() *ConfigureManager {
	except.ASSERT(configureManager != nil)
	return configureManager
}

func (p *ConfigureManager) GetConfigItem(name string) *Item {
	itemobj, ok := p.cache.Load(name)
	except.ASSERT(ok)

	return itemobj.(*Item)
}

func (p *ConfigureManager) GetConfigItemFromDb(name string) *Item {
	config := Configure{}
	db := dbo.DboInstance().DB()
	err := db.Where("param_name=?", name).Find(&config).Error
	if err != nil {
		Log.Error("db error = %v", err)
		panic(err)
	}

	return &Item{Value: config.ParamValue, Disabled: config.Disabled}
}

func (p *ConfigureManager) SetConfigItemWithEncryption(name string, value string) error {

	encrypted, err := encrypt.InternalEncryptStr(value)
	if err != nil {
		Log.Errorf("value=%v, err=%v", value, err)
		return err
	}

	return p.SetConfigItem(name, encrypted)
}

func (p *ConfigureManager) SetConfigItem(name string, value string) error {

	config := &Configure{
		ParamName:  name,
		ParamValue: value,
	}

	db := dbo.DboInstance().DB()
	err := db.Where("param_name=?", name).Update("param_value", value).Error
	if err != nil {
		Log.Errorf("db error = %v", err)
		return err
	}

	p.cache.Store(name, config)

	return nil
}

func GetItem(name string) *Item {
	return ConfigureManagerInstance().GetConfigItem(name)
}

func SetItem(name string, value string) error {
	return ConfigureManagerInstance().SetConfigItem(name, value)
}

func GetItemDirect(name string) *Item {
	return ConfigureManagerInstance().GetConfigItemFromDb(name)
}

func SetItemWithEncryption(name string, value string) error {
	return ConfigureManagerInstance().SetConfigItemWithEncryption(name, value)
}
