package dbo

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/nece099/base/dbo/dblogger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	DB_TYPE_MYSQL  = "mysql"
	DB_TYPE_SQLITE = "sqlite"
)

type DboConfig struct {
	DbType      string
	URL         string
	IdleSize    int
	MaxSize     int
	MaxLifeTime int64
	SqlDebug    int
	AutoMigrate bool
}

type Dbo struct {
	config *gorm.Config
	db     *gorm.DB
	models []interface{}
}

var dbo *Dbo = &Dbo{}

func openDb(c *DboConfig) (*gorm.DB, error) {
	dbtype := c.DbType
	if len(dbtype) == 0 {
		dbtype = DB_TYPE_MYSQL
	}

	if dbtype == DB_TYPE_MYSQL {
		db, err := gorm.Open(mysql.Open(c.URL), dbo.config)
		return db, err
	} else if dbtype == DB_TYPE_SQLITE {
		db, err := gorm.Open(sqlite.Open(c.URL), dbo.config)
		return db, err
	} else {
		return nil, errors.New(fmt.Sprintf("unsupported db type:%v", dbtype))
	}
}

func DboInit(configs []*DboConfig) {

	if len(configs) == 0 {
		panic("no db config found!")
	}
	c := configs[0]

	logLv := glogger.Silent
	if c.SqlDebug == 1 {
		logLv = glogger.Info
	}

	dbo.config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: &dblogger.DbLogger{
			LogLevel: logLv,
		},
	}

	db, err := openDb(c)
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

	// 设置字符编码
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	for _, m := range dbo.models {
		if !db.Migrator().HasTable(m) {
			err := db.Migrator().CreateTable(m)
			if err != nil {
				Log.Errorf("m = %v, err = %v", reflect.TypeOf(m), err)
			}
		}
	}

	if c.AutoMigrate {
		db.AutoMigrate(dbo.models...)
	}

	dbo.db = db
}

func RegisterModels(models ...interface{}) {
	dbo.models = models
}

func DboInstance() *Dbo {
	ASSERT(dbo != nil)
	return dbo
}

func (s *Dbo) DB() *gorm.DB {
	db := s.db
	sessdb := db.Session(&gorm.Session{})
	return sessdb
}
