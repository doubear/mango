package mango

import (
	"os"

	"github.com/go-mango/mango/logger"
)

var mlog *logger.Logger

func init() {
	mlog = logger.NewLogger()
}

//Debug log debug message with given format.
func Debug(s string, a ...interface{}) {
	mlog.Debug(s, a...)
}

//Info log info message with given format.
func Info(s string, a ...interface{}) {
	mlog.Info(s, a...)
}

//Warn log warn message with given format.
func Warn(s string, a ...interface{}) {
	mlog.Warn(s, a...)
}

//Fatal log fatal message with given format.
func Fatal(s string, a ...interface{}) {
	mlog.Fatal(s, a...)
}

//SetLevel set logger level.
func SetLevel(l int) {
	mlog.SetLevel(l)
}

//SetOutput set logger output.
func SetOutput(w *os.File) {
	mlog.SetOutput(w)
}

//SetPrefix set log prefix.
func SetPrefix(s string) {
	mlog.SetPrefix(s)
}
