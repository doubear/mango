package mango

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

type message struct {
	id      int64
	level   int
	content string
	time    time.Time
	logger  *Logger
}

func (this *message) color() int {
	return colors[this.level]
}

func (this *message) fillColor(s string, c int) string {
	if this.logger.flag&FlagColor != 0 {
		if c == 0 {
			c = this.color()
		}

		return "\033[" + strconv.Itoa(c) + "m" + s + "\033[0m"
	}

	return s
}

func (this *message) label() string {
	return this.fillColor(levels[this.level], 0)
}

//Stringer convert message bag into string.
//string format: [datetime] [id] [level] [content]
func (this *message) toString() string {
	var stack []string

	if this.logger.flag&dateMask != 0 {
		fTime := DateFormat(this.time, this.logger.dateFormat())
		stack = append(stack, fTime)
	}

	if this.logger.flag&FlagID != 0 {
		stack = append(stack, "id:"+int2Str(this.id, 0))
	}

	if this.logger.flag&FlagLevel != 0 {
		stack = append(stack, this.label())
	}

	stack = append(stack, this.logger.prefix+this.content)

	return "\r" + strings.Join(stack, " ") + "\n"
}

//Consume write message to endpoint.
func (this *message) consume() {
	this.logger.out.Write([]byte(this.toString()))
}

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
func (this *Logger) SetPrefix(s string) {
	this.prefix = s
}

//SetOutput reset logs endpoint.
func (this *Logger) SetOutput(o *os.File) {
	this.out = o
}

//SetLevel resets loggable min level.
func (this *Logger) SetLevel(l int) {
	this.logable = l
}

//DateFormat returns datetime format struct by flags.
func (this *Logger) dateFormat() string {
	var date []string
	var time []string

	if this.flag&FlagYear != 0 {
		date = append(date, "Y")
	}

	if this.flag&FlagMonth != 0 {
		date = append(date, "m")
	}

	if this.flag&FlagDay != 0 {
		date = append(date, "d")
	}

	if this.flag&FlagHour != 0 {
		time = append(time, "H")
	}

	if this.flag&FlagMinute != 0 {
		time = append(time, "i")
	}

	if this.flag&FlagSecond != 0 {
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
func (this *Logger) newMessage(l int, s string) *message {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.count++

	return &message{
		this.count,
		l,
		s,
		time.Now(),
		this,
	}
}

func (this *Logger) print(l int, s string) {
	if l >= this.logable {
		this.newMessage(l, s).consume()
	}
}

func (this *Logger) printf(l int, s string, a ...interface{}) {
	if l >= this.logable {
		this.newMessage(l, fmt.Sprintf(s, a...)).consume()
	}
}

//Debug writes a debug level log.
func (this *Logger) Debug(s string) {
	this.print(LogDebug, s)
}

func (this *Logger) Debugf(s string, a ...interface{}) {
	this.printf(LogDebug, s, a...)
}

//Info writes a info level log.
func (this *Logger) Info(s string) {
	this.print(LogInfo, s)
}

func (this *Logger) Infof(s string, a ...interface{}) {
	this.printf(LogInfo, s, a...)
}

//Warn writes a warn level log.
func (this *Logger) Warn(s string) {
	this.print(LogWarn, s)
}

func (this *Logger) Warnf(s string, a ...interface{}) {
	this.printf(LogWarn, s, a...)
}

//Fatal writes a fatal level log.
func (this *Logger) Fatal(s string) {
	this.print(LogFatal, s)
}

func (this *Logger) Fatalf(s string, a ...interface{}) {
	this.printf(LogFatal, s, a...)
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

var defaultLogger = NewLogger()
