package dbo

import (
	"os"
	"reflect"
	"time"

	"github.com/nece099/base/dbo/dblogger"
	"github.com/nece099/base/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DboConfig struct {
	URL         string
	IdleSize    int
	MaxSize     int
	MaxLifeTime int64
	SqlDebug    int
}

type Dbo struct {
	config *gorm.Config
	db     *gorm.DB
	models []interface{}
}

var dbo *Dbo = nil

func DboInit(configs []*DboConfig) {

	if len(configs) == 0 {
		panic("no db config found!")
	}
	c := configs[0]

	logLv := glogger.Silent
	if c.SqlDebug == 1 {
		logLv = glogger.Info
	}

	dbo = &Dbo{
		config: &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			},
			Logger: &dblogger.DbLogger{
				LogLevel: logLv,
			},
		},
	}

	db, err := gorm.Open(mysql.Open(c.URL), dbo.config)
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}

	sdb, err := db.DB()
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}

	sdb.SetMaxIdleConns(c.IdleSize)
	sdb.SetMaxOpenConns(c.MaxSize)
	sdb.SetConnMaxLifetime(time.Duration(c.MaxLifeTime) * time.Second)

	for _, m := range dbo.models {
		if !db.Migrator().HasTable(m) {
			err := db.Migrator().CreateTable(m)
			if err != nil {
				Log.Errorf("m = %v, err = %v", reflect.TypeOf(m), err)
			}
		}
	}
	db.AutoMigrate(dbo.models...)

	dbo.db = db
}

func RegisterModels(models ...interface{}) {
	dbo.models = models
}

func DboInstance() *Dbo {
	utils.ASSERT(dbo != nil)
	return dbo
}

func (s *Dbo) DB() *gorm.DB {
	return dbo.db
}
