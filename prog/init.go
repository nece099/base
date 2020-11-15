package prog

import "github.com/zen099/onetube/server/base/logger"

var Log *logger.Logger = nil

func init() {
	Log = logger.Log
}
