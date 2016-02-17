package common

import (
	log "github.com/Sirupsen/logrus"
	"runtime"
)

var Log *log.Entry

func init() {
	Log = log.NewEntry(log.StandardLogger())
}

//recover after panic, log traceback
func LogErr() {
	if err := recover(); err != nil {
		trace := make([]byte, 10024)
		count := runtime.Stack(trace, true)
		Log.WithField("error", err).Printf("[error recovered]")
		Log.Printf("Stack trace, lines %d  trace: %s", count, string(trace))
	}
}
