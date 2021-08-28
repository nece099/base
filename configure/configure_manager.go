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
		c := &Configure{
			ParamValue: v,
			Disabled:   false,
		}

		configureManager.cache.Store(k, c)
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
		configureManager.cache.Store(c.ParamName, c)
	}
}

func ConfigureManagerInstance() *ConfigureManager {
	except.ASSERT(configureManager != nil)
	return configureManager
}

func (p *ConfigureManager) GetConfigItem(name string) *Configure {
	itemobj, ok := p.cache.Load(name)
	except.ASSERT(ok)

	return itemobj.(*Configure)
}

func (p *ConfigureManager) GetConfigItemFromDB(name string) *Configure {
	config := &Configure{}
	db := dbo.DboInstance().DB()
	err := db.Model(&Configure{}).Where("param_name=?", name).First(config).Error
	if err != nil {
		Log.Error("db error = %v", err)
		panic(err)
	}

	return config
}

func (p *ConfigureManager) SetConfigItem(name string, value string) error {

	config := &Configure{
		ParamName: name,
	}

	db := dbo.DboInstance().DB()
	err := db.Model(&Configure{}).First(config).Error
	if err != nil {
		Log.Error(err)
		return err
	}

	if config.Encrypted {
		val, err := encrypt.InternalEncryptStr(value)
		if err != nil {
			Log.Error(err)
			return err
		}

		config.ParamValue = val
	} else {
		config.ParamValue = value
	}

	err = db.Model(&Configure{}).
		Where("param_name=?", name).
		Update("param_value", config.ParamValue).Error
	if err != nil {
		Log.Errorf("db error = %v", err)
		return err
	}

	p.cache.Store(name, config)

	return nil
}

func (p *ConfigureManager) ConfigureList() []*Configure {
	cfgs := []*Configure{}
	db := dbo.DboInstance().DB()
	err := db.Model(&Configure{}).
		Where("hidden=0").
		Order("id").
		Find(&cfgs).Error
	if err != nil {
		Log.Error("db error = %v", err)
		panic(err)
	}

	return cfgs
}

func GetItem(name string) *Configure {
	return ConfigureManagerInstance().GetConfigItem(name)
}

func SetItem(name string, value string) error {
	return ConfigureManagerInstance().SetConfigItem(name, value)
}

func GetItemDirect(name string) *Configure {
	return ConfigureManagerInstance().GetConfigItemFromDB(name)
}

func ConfigureList() []*Configure {
	return ConfigureManagerInstance().ConfigureList()
}
