package middlewares

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/contracts"
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
func Throttle(qps int) contracts.ThenableFunc {

	if qps < 0 {
		logy.Std().Error("ThrottleOption QPS must larger than 0")
	}

	if qps == 0 {
		qps = 15
		logy.Std().Warn("ThrottleOption uses default rate 15 req/s")
	}

	var hashmap = make(map[string]*throttle) //summary & times

	return func(ctx contracts.ThenableContext) {
		label := ctx.Request().IP() + ctx.Request().URI()
		barr := sha1.Sum([]byte(label))
		sum := hex.EncodeToString(barr[:])

		if t, ok := hashmap[sum]; ok {
			t.m.Lock()
			defer t.m.Unlock()

			if time.Since(t.t) <= 1*time.Second {

				if t.c >= qps {
					ctx.Response().SetStatus(http.StatusTooManyRequests)
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
