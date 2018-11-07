package logy

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

//Logger flags
const (
	FlagLevel = 1 << iota
	FlagColor
	FlagID
	FlagDate
)

//log levels
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
)

const (
	fBlack = 30 + iota
	fRed
	fGreen
	fYellow
	fBlue
	fPurple
	fCyan
	fWhite
	fNone = 0
)

var (
	levels = []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
	}

	colors = map[int]int{
		LogDebug: fNone,
		LogInfo:  fGreen,
		LogWarn:  fYellow,
		LogError: fRed,
	}

	flag = FlagColor | FlagLevel | FlagDate
)

//Logger 日志记录器结构
type Logger struct {
	sync.Mutex
	flag       int
	count      int64
	out        io.Writer
	prefix     string
	logable    int
	dateFormat string
}

//SetPrefix 设置日志消息前缀
func (l *Logger) SetPrefix(s string) {
	l.prefix = s
}

//SetOutput 指定日志输出器
func (l *Logger) SetOutput(o io.Writer) {
	l.out = o
}

//SetLevel 指定日志的记录起始等级
func (l *Logger) SetLevel(lv int) {
	l.logable = lv
}

//SetDateFormat 设置日期格式
//参考：https://golang.org/src/time/format.go
func (l *Logger) SetDateFormat(f string) {
	l.dateFormat = f
}

//createMessage 创建一个消息
func (l *Logger) createMessage(lv int, s string) {
	l.Lock()
	defer l.Unlock()

	l.count++

	m := &message{
		l.count,
		lv,
		s,
		time.Now(),
		l,
	}

	m.consume()
}

func (l *Logger) print(lv int, s interface{}, a ...interface{}) {
	if lv >= l.logable {
		_, file, line, _ := runtime.Caller(3)

		f := fmt.Sprintf("%s:%d %s", chopPath(file), line, formatString(s))

		l.createMessage(lv, fmt.Sprintf(f, a...))
	}
}

//D 通过 a 格式化 s 并输出为 debug 等级的日志
func (l *Logger) D(s interface{}, a ...interface{}) {
	l.print(LogDebug, s, a...)
}

//I 通过 a 格式化 s 并输出为 debug 等级的日志
func (l *Logger) I(s interface{}, a ...interface{}) {
	l.print(LogInfo, s, a...)
}

//W 通过 a 格式化 s 并输出为 debug 等级的日志
func (l *Logger) W(s interface{}, a ...interface{}) {
	l.print(LogWarn, s, a...)
}

//E 通过 a 格式化 s 并输出为 debug 等级的日志
func (l *Logger) E(s interface{}, a ...interface{}) {
	l.print(LogError, s, a...)
	os.Exit(0)
}

func formatString(s interface{}) string {
	f := "%v"

	switch s.(type) {
	case error:
		f = s.(error).Error()
	case fmt.Stringer:
		f = s.(fmt.Stringer).String()
	case string:
		f = s.(string)
	}

	return f
}

func chopPath(s string) string {
	_, f := path.Split(s)

	return f
}

//New create logger instance.
func New() *Logger {
	return &Logger{
		sync.Mutex{},
		flag,
		0,
		os.Stdout,
		"",
		LogDebug,
		"15:04:05",
	}
}
