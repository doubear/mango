package logy

import "io"

var l *Logger

func init() {
	l = New()
}

//D 写入debug级日志信息
func D(s interface{}, a ...interface{}) {
	l.D(s, a...)
}

//I 写入 info 级日志信息
func I(s interface{}, a ...interface{}) {
	l.I(s, a...)
}

//W 写入 warn 级日志信息
func W(s interface{}, a ...interface{}) {
	l.W(s, a...)
}

//E 写入 error 级日志信息
func E(s interface{}, a ...interface{}) {
	l.E(s, a...)
}

//SetLevel 设置需要记录的最低日志等级
func SetLevel(lv int) {
	l.SetLevel(lv)
}

//SetOutput 设置日志输出文件
func SetOutput(w io.Writer) {
	l.SetOutput(w)
}

//SetPrefix 设置日志信息前缀
func SetPrefix(s string) {
	l.SetPrefix(s)
}

//SetDateFormat 设置日志创建时间记录格式
func SetDateFormat(f string) {
	l.SetDateFormat(f)
}
