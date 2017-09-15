package mango

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/go-mango/logy"
)

type throttle struct {
	c int
	t time.Time
	m sync.Mutex
}

func (t *throttle) reset() {
	t.c = 0
	t.t = time.Now()
}

//Throttle controls how much request frequency
//from remote client is allowed.
func Throttle(qps int) MiddleFunc {

	if qps < 0 {
		logy.E("ThrottleOption QPS must larger than 0")
	}

	if qps == 0 {
		qps = 15
		logy.W("ThrottleOption uses default rate 15 req/s")
	}

	var hashmap = make(map[string]*throttle) //summary & times

	return func(ctx *Context) {
		label := ctx.ClientIP() + ctx.R.RequestURI
		barr := sha1.Sum([]byte(label))
		sum := hex.EncodeToString(barr[:])

		if t, ok := hashmap[sum]; ok {
			t.m.Lock()
			defer t.m.Unlock()

			if time.Since(t.t) <= 1*time.Second {

				if t.c >= qps {
					ctx.W.SetStatus(http.StatusTooManyRequests)
					return
				}

			} else {
				t.reset()
			}

			t.c++
		} else {
			hashmap[sum] = &throttle{1, time.Now(), sync.Mutex{}}
		}

		ctx.Next()
	}
}
