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

const (
	CONN_NORMAL = "normal" // 一般db
	CONN_REP    = "rep"    // 复制库 只读, 读写分离
	CONN_MEM    = "mem"    // 内存表
)

type DboConfig struct {
	DbType      string
	URL         string
	IdleSize    int
	MaxSize     int
	MaxLifeTime int64
	SqlDebug    int
	AutoMigrate bool
	ConnType    string
}

type Dbo struct {
	config *gorm.Config
	db     *gorm.DB
	repDb  *gorm.DB
	memDb  *gorm.DB
	models []interface{}
}

var dbo *Dbo = &Dbo{}

func openDb(c *DboConfig) (*gorm.DB, error) {
	dbtype := c.DbType
	if len(dbtype) == 0 {
		dbtype = DB_TYPE_MYSQL
	}

	connType := c.ConnType
	if len(connType) == 0 {
		connType = CONN_NORMAL
	}

	if dbtype == DB_TYPE_MYSQL {
		if len(c.URL) == 0 { // 没有url则忽略
			return nil, nil
		}

		db, err := gorm.Open(mysql.Open(c.URL), dbo.config)
		// 设置字符编码
		if connType == CONN_MEM {
			db = db.Set("gorm:table_options", "ENGINE=MEMORY CHARSET=utf8mb4")
		} else {
			db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
		}

		sdb, err := db.DB()
		if err != nil {
			Log.Error(err)
			os.Exit(-1)
		}

		sdb.SetMaxIdleConns(c.IdleSize)
		sdb.SetMaxOpenConns(c.MaxSize)
		sdb.SetConnMaxLifetime(time.Duration(c.MaxLifeTime) * time.Second)

		// 一般db和memdb才做自动建表
		if connType == CONN_NORMAL || connType == CONN_MEM {
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
		}

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

	// normal db
	db, err := openDb(c)
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}

	dbo.db = db

	// rep db
	if len(configs) >= 2 {
		c1 := configs[1]

		if c1.ConnType != CONN_REP {
			Log.Error("the 2nd config have to be rep db")
			os.Exit(-1)
		}

		// rep db
		repdb, err := openDb(c1)
		if err != nil {
			Log.Error(err)
			os.Exit(-1)
		}

		dbo.repDb = repdb
	}

	// mem db
	if len(configs) >= 3 {
		c2 := configs[2]

		if c2.ConnType != CONN_MEM {
			Log.Error("the 3th config have to be mem db")
			os.Exit(-1)
		}

		// rep db
		memdb, err := openDb(c2)
		if err != nil {
			Log.Error(err)
			os.Exit(-1)
		}

		dbo.memDb = memdb
	}
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

func (s *Dbo) RepDB() *gorm.DB {
	db := s.repDb
	if db != nil {
		sessdb := db.Session(&gorm.Session{})
		return sessdb
	}

	return nil
}

func (s *Dbo) MemDB() *gorm.DB {
	db := s.memDb
	if db != nil {
		sessdb := db.Session(&gorm.Session{})
		return sessdb
	}

	return nil
}
