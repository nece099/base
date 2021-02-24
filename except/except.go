package except

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/nece099/base/utils"
)

func CatchPanic() {
	if err := recover(); err != nil {
		Log.Errorf("panic !!! err = %v ", err)
	}
}

func CatchPanicWarning() {
	if err := recover(); err != nil {
		Log.Warnf("panic !!! err = %v ", err)
	}
}

func CatchException() {
	if err := recover(); err != nil {
		fullPath, _ := exec.LookPath(os.Args[0])
		fname := filepath.Base(fullPath)

		datestr := utils.NowDateStr()
		outstr := fmt.Sprintf("\n======\n[%v] err=%v, stack=%v\n======\n", time.Now(), err, string(debug.Stack()))
		filename := "./log/panic_" + fname + datestr + ".log"
		f, err2 := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		ASSERT(err2 == nil)
		defer f.Close()
		f.WriteString(outstr)

		Log.Errorf("err = %v ", err)
	}
}
