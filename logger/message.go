package logger

import (
	"strconv"
	"strings"
	"time"
)

type message struct {
	id      int64
	level   int
	content string
	time    time.Time
	logger  *Logger
}

func (msg *message) color() int {
	return colors[msg.level]
}

func (msg *message) fillColor(s string, c int) string {
	if msg.logger.flag&FlagColor != 0 {
		if c == 0 {
			c = msg.color()
		}

		return "\033[" + strconv.Itoa(c) + "m" + s + "\033[0m"
	}

	return s
}

func (msg *message) label() string {
	return msg.fillColor(levels[msg.level], 0)
}

//Stringer convert message bag into string.
//string format: [datetime] [id] [level] [content]
func (msg *message) toString() string {
	var stack []string

	if msg.logger.flag&dateMask != 0 {
		fTime := DateFormat(msg.time, msg.logger.dateFormat())
		stack = append(stack, fTime)
	}

	if msg.logger.flag&FlagID != 0 {
		stack = append(stack, "id:"+int2Str(msg.id, 0))
	}

	if msg.logger.flag&FlagLevel != 0 {
		stack = append(stack, msg.label())
	}

	stack = append(stack, msg.logger.prefix+msg.content)

	return "\r" + strings.Join(stack, " ") + "\n"
}

//Consume write message to endpoint.
func (msg *message) consume() {
	msg.logger.out.Write([]byte(msg.toString()))
}
