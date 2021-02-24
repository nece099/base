package dbo

import (
	"github.com/nece099/base/logger"
)

var Log *logger.Logger = nil

func init() {
	Log = logger.Log
}

func ASSERT(b bool) {
	if !b {
		panic("Assert failed")
	}
}
