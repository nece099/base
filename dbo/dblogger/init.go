package gormlogger

import (
	"github.com/maneki001/tgflow/server/base/logger"
)

var Log *logger.Logger = nil

func init() {
	Log = logger.Log
}
