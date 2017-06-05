package mango

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//DateFormat format time.Time into given format.
//Allowed labels are list below:
//Y: year
//m: Month
//d: date
//H: hour
//i: minute
//s: second
func DateFormat(t time.Time, f string) string {
	if f == "" {
		f = "Y-m-d H:i:s"
	}

	y, m, d := t.Date()
	h, i, s := t.Clock()

	f = strings.Replace(f, "Y", int2Str(int64(y), 0), -1)
	f = strings.Replace(f, "m", int2Str(int64(int(m)), 2), -1)
	f = strings.Replace(f, "d", int2Str(int64(d), 2), -1)
	f = strings.Replace(f, "H", int2Str(int64(h), 2), -1)
	f = strings.Replace(f, "i", int2Str(int64(i), 2), -1)
	f = strings.Replace(f, "s", int2Str(int64(s), 2), -1)

	return f
}

//Int2Str convert int to string with prefix padding.
func int2Str(i int64, w int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(w)+"d", i)
}

var traces = map[string]time.Time{}

func StartTrace(id string) {
	traces[id] = time.Now()
}

func StopTrace(id string) {
	e := time.Since(traces[id])
	delete(traces, id)
	defaultLogger.Debugf("duration %s spent %dns", id, e.Nanoseconds())
}

//NumericTimeSmartFormat format the given nanosecond.
func NumericTimeSmartFormat(t int64) string {
	var f float64
	switch {
	case t > 1000000000: //s
		f = float64(t) / 1000000000
		return fmt.Sprintf("%.2fs", f)
	case t > 1000000: //ms
		f = float64(t) / 1000000
		return fmt.Sprintf("%.2fms", f)
	case t > 1000: //us
		f = float64(t) / 1000
		return fmt.Sprintf("%.2fμs", f)
	}

	return fmt.Sprintf("%dns", t)
}
