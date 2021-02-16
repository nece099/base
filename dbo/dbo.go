package dbo

import (
	"os"
	"reflect"
	"time"

	"github.com/nece099/base/dbo/dblogger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// const (
// 	IDLE_CONN          = 1
// 	MAX_OPEN_CONN      = 100
// 	CONN_MAX_LIFE_TIME = 120
// )

type DboConfig struct {
	DSN                   string
	IdleConnection        int
	MaxOpenConnection     int
	ConnectionMaxLifeTime int
}

type Dbo struct {
	config *gorm.Config
	db     *gorm.DB
	models []interface{}
}

var dbologger = glogger.New()

var dbo *Dbo = &Dbo{
	config: &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: &dblogger.DbLogger{
			LogLevel: glogger.Silent,
		},
	},
}

func DboInit(c *DboConfig) {
	db, err := gorm.Open(mysql.Open(c.DSN), dbo.config)
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}

	sdb, err := db.DB()
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}

	sdb.SetMaxIdleConns(c.IdleConnection)
	sdb.SetMaxOpenConns(c.MaxOpenConnection)
	sdb.SetConnMaxLifetime(time.Duration(c.ConnectionMaxLifeTime) * time.Second)

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

func DboRegisterModels(models ...interface{}) {
	dbo.models = models
}

func DboInstance() *Dbo {
	return dbo
}

func (s *Dbo) DB() *gorm.DB {
	return dbo.db
}
