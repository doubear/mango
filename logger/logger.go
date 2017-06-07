package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

//Logger flags
const (
	FlagYear = 1 << iota
	FlagMonth
	FlagDay
	FlagHour
	FlagMinute
	FlagSecond
	FlagLevel
	FlagColor
	FlagID
	dateMask = FlagYear | FlagMonth | FlagDay | FlagHour | FlagMinute | FlagSecond
)

//log levels
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogFatal
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
		"FATAL",
	}

	colors = map[int]int{
		LogDebug: fNone,
		LogInfo:  fGreen,
		LogWarn:  fYellow,
		LogFatal: fRed,
	}

	flag = FlagColor | FlagLevel | FlagHour | FlagMinute | FlagSecond
)

//Logger log system main entry
type Logger struct {
	mu      sync.Mutex
	flag    int
	count   int64
	out     io.Writer
	prefix  string
	logable int
}

//SetPrefix reset prefix of every logs.
func (log *Logger) SetPrefix(s string) {
	log.prefix = s
}

//SetOutput reset logs endpoint.
func (log *Logger) SetOutput(o *os.File) {
	log.out = o
}

//SetLevel resets loggable min level.
func (log *Logger) SetLevel(l int) {
	log.logable = l
}

//DateFormat returns datetime format struct by flags.
func (log *Logger) dateFormat() string {
	var date []string
	var time []string

	if log.flag&FlagYear != 0 {
		date = append(date, "Y")
	}

	if log.flag&FlagMonth != 0 {
		date = append(date, "m")
	}

	if log.flag&FlagDay != 0 {
		date = append(date, "d")
	}

	if log.flag&FlagHour != 0 {
		time = append(time, "H")
	}

	if log.flag&FlagMinute != 0 {
		time = append(time, "i")
	}

	if log.flag&FlagSecond != 0 {
		time = append(time, "s")
	}

	var datetime string

	if len(time) > 0 {
		datetime = strings.Join(time, ":")
	}

	if len(date) > 0 {
		datetime = strings.Join(date, "-") + " " + datetime
	}

	return strings.TrimSpace(datetime)
}

//newMessage create new log message.
func (log *Logger) newMessage(l int, s string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.count++

	m := &message{
		log.count,
		l,
		s,
		time.Now(),
		log,
	}

	m.consume()
}

func (log *Logger) print(l int, s string) {
	if l >= log.logable {
		log.newMessage(l, s)
	}
}

func (log *Logger) printf(l int, s string, a ...interface{}) {
	if l >= log.logable {
		log.newMessage(l, fmt.Sprintf(s, a...))
	}
}

//Debug writes a debug level log.
func (log *Logger) Debug(s string) {
	log.print(LogDebug, s)
}

//Debugf format s with a.
func (log *Logger) Debugf(s string, a ...interface{}) {
	log.printf(LogDebug, s, a...)
}

//Info writes a info level log.
func (log *Logger) Info(s string) {
	log.print(LogInfo, s)
}

//Infof format s with a.
func (log *Logger) Infof(s string, a ...interface{}) {
	log.printf(LogInfo, s, a...)
}

//Warn writes a warn level log.
func (log *Logger) Warn(s string) {
	log.print(LogWarn, s)
}

//Warnf format s with a.
func (log *Logger) Warnf(s string, a ...interface{}) {
	log.printf(LogWarn, s, a...)
}

//Fatal writes a fatal level log.
func (log *Logger) Fatal(s string) {
	log.print(LogFatal, s)
}

//Fatalf format s with a.
func (log *Logger) Fatalf(s string, a ...interface{}) {
	log.printf(LogFatal, s, a...)
}

//NewLogger create logger instance.
func NewLogger() *Logger {
	return &Logger{
		sync.Mutex{},
		flag,
		0,
		os.Stdout,
		"",
		LogDebug,
	}
}
