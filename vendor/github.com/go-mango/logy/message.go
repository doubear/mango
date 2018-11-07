package logy

import (
	"fmt"
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

func (m *message) color() int {
	return colors[m.level]
}

func (m *message) fillColor(s string, c int) string {
	if m.logger.flag&FlagColor != 0 {
		if c == 0 {
			c = m.color()
		}

		return "\033[" + strconv.Itoa(c) + "m" + s + "\033[0m"
	}

	return s
}

func (m *message) label() string {
	return m.fillColor(levels[m.level], 0)
}

//日志完整格式： [id] [datetime] [level] [content]
func (m *message) toString() string {
	var stack []string

	if m.logger.flag&FlagID != 0 {
		stack = append(stack, fmt.Sprintf("id: %d", m.id))
	}

	if m.logger.flag&FlagDate != 0 {
		stack = append(stack, m.time.Format(m.logger.dateFormat))
	}

	if m.logger.flag&FlagLevel != 0 {
		stack = append(stack, m.label())
	}

	stack = append(stack, m.logger.prefix, m.content)

	return "\r" + strings.Join(stack, " ") + "\n"
}

func (m *message) consume() {
	m.logger.out.Write([]byte(m.toString()))
}
