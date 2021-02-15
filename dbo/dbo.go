package dbo

import (
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Dbo struct {
	config *gorm.Config
	db     *gorm.DB
	models []interface{}
}

var dbo *Dbo = &Dbo{
	config: &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: &gormlogger.GormLogger{
			LogLevel: glogger.Silent,
		},
	},
}
