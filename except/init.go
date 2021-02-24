package except

import (
	"github.com/nece099/base/logger"
)

var Log *logger.Logger = nil

func init() {
	Log = logger.Log
}
