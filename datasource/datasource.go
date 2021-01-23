package datasource

import (
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "github.com/nece099/base/logger"
)

type TSQLLogger struct{}

func (slog TSQLLogger) Print(values ...interface{}) {
	vals := gorm.LogFormatter(values...)
	Log.SqlDebug(vals...)
}

type DataSourceManager struct {
	configs      []*DataSource
	dbs          []*gorm.DB
	masterDB     *gorm.DB
	slaveDB      *gorm.DB
	memoryDB     *gorm.DB
	models       []interface{}
	memoryModels []interface{}
}

var datasourceManager = &DataSourceManager{}

type DataSource struct {
	URL         string
	IdleSize    int
	MaxSize     int
	MaxLifeTime int64
	SqlDebug    int
	Memory      bool
}

func DataSourceInit(configs []*DataSource) {
	datasourceManager.configs = configs

	if len(configs) == 0 {
		Log.Warnf("db config not found...")
		return
	}

	for idx, conf := range configs {

		var db *gorm.DB
		if conf.Memory {
			db = datasourceManager.openMemDBConn(conf)
		} else {
			db = datasourceManager.openDBConn(conf)
		}

		datasourceManager.dbs = append(datasourceManager.dbs, db)

		if idx == 0 {
			datasourceManager.masterDB = db
			datasourceManager.slaveDB = db
			datasourceManager.memoryDB = db
		}

		if idx == 1 { // 如果有配置从库
			datasourceManager.slaveDB = db
		}

		if datasourceManager.memoryDB == nil && conf.Memory { // 只支持一个内存库
			datasourceManager.memoryDB = db
		}
	}
}

func DataSourceRegisterModels(models ...interface{}) {
	datasourceManager.models = models
}

func DataSourceRegisterMemoryModels(models ...interface{}) {
	datasourceManager.memoryModels = models
}

func (d *DataSourceManager) openDBConn(ds *DataSource) *gorm.DB {
	db, err := gorm.Open("mysql", ds.URL)
	if err != nil {
		Log.Errorf("connect to mysql failed, err = %v", err)
		return nil
	}

	db.DB().SetMaxIdleConns(ds.IdleSize)
	db.DB().SetMaxOpenConns(ds.MaxSize)
	db.DB().SetConnMaxLifetime(time.Duration(ds.MaxLifeTime) * time.Second)

	// 设置字符编码
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	db.SingularTable(true)
	if ds.SqlDebug == 1 {
		db.LogMode(true)
		db.SetLogger(TSQLLogger{})
	}

	for _, m := range d.models {
		if !db.HasTable(m) {
			// Log.Debugf("m = %v", reflect.TypeOf(m))
			err := db.CreateTable(m).Error
			if err != nil {
				Log.Errorf("m = %v, err = %v", reflect.TypeOf(m), err)
			}
		}
	}
	db.AutoMigrate(d.models...)

	return db
}

func (d *DataSourceManager) openMemDBConn(ds *DataSource) *gorm.DB {
	db, err := gorm.Open("mysql", ds.URL)
	if err != nil {
		Log.Errorf("connect to mysql failed, err = %v", err)
		return nil
	}

	db.DB().SetMaxIdleConns(ds.IdleSize)
	db.DB().SetMaxOpenConns(ds.MaxSize)
	db.DB().SetConnMaxLifetime(time.Duration(ds.MaxLifeTime) * time.Second)

	// 设置字符编码
	db = db.Set("gorm:table_options", "ENGINE=MEMORY CHARSET=utf8mb4")
	db.SingularTable(true)
	if ds.SqlDebug == 1 {
		db.LogMode(true)
		db.SetLogger(TSQLLogger{})
	}

	for _, m := range d.memoryModels {
		if !db.HasTable(m) {
			err := db.CreateTable(m).Error
			if err != nil {
				Log.Error(err)
			}
		}
	}
	db.AutoMigrate(d.memoryModels...)

	return db
}

func (d *DataSourceManager) Master() *gorm.DB {
	return d.masterDB
}

func (d *DataSourceManager) Slave() *gorm.DB {
	return d.slaveDB
}

func (d *DataSourceManager) Memory() *gorm.DB {
	return d.memoryDB
}

func DataSourceInstance() *DataSourceManager {
	return datasourceManager
}
