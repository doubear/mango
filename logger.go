package mango

import (
	"os"

	"github.com/go-mango/mango/logger"
)

var mlog *logger.Logger

func init() {
	mlog = logger.NewLogger()
}

//Debugf log debug message with given format.
func Debugf(s string, a ...interface{}) {
	mlog.Debugf(s, a...)
}

//Infof log info message with given format.
func Infof(s string, a ...interface{}) {
	mlog.Infof(s, a...)
}

//Warnf log warn message with given format.
func Warnf(s string, a ...interface{}) {
	mlog.Warnf(s, a...)
}

//Fatalf log fatal message with given format.
func Fatalf(s string, a ...interface{}) {
	mlog.Fatalf(s, a...)
}

//Debug log debug message.
func Debug(s string) {
	Debugf(s)
}

//Info log Info message.
func Info(s string) {
	Infof(s)
}

//Warn log Warn message.
func Warn(s string) {
	Warnf(s)
}

//Fatal log Fatal message.
func Fatal(s string) {
	Fatalf(s)
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
