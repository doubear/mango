package mango

import (
	"fmt"
	"net/http"
	"time"
)

//Record log incoming requests to console.
func Record() MiddleFunc {
	return func(ctx *Context) {
		st := time.Now()
		ctx.Next()
		dur := NumericTimeSmartFormat(time.Since(st).Nanoseconds())

		ctx.Logger.Infof(
			"%s %s %s %d %dbytes %s",
			ctx.ClientIP(),
			ctx.R.Method,
			ctx.R.RequestURI,
			ctx.W.Status(),
			ctx.W.Size(),
			dur,
		)
	}
}

//Recovery record every panic.
func Recovery() MiddleFunc {
	return func(ctx *Context) {
		defer func() {
			if v := recover(); v != nil {
				var err error
				switch v := v.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", v)
				}

				ctx.Logger.Warn(err.Error())

				ctx.W.SetStatus(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}

//Cors addtional CORS middleware
func Cors() MiddleFunc {
	return func(ctx *Context) {
		if ctx.R.Method == "OPTIONS" {
			ctx.W.Header().Add("Access-Control-Allow-Origin", "*")
			ctx.W.Header().Add("Access-Control-Allow-Methods", "*")
			ctx.W.Header().Add("Access-Control-Allow-Headers", "*")
			ctx.W.Header().Add("Access-Control-Max-Age", "86400")
			ctx.W.SetStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}

//Throttle controls how much request frequency
//from remote client is allowed.
func Throttle() MiddleFunc {
	return func(ctx *Context) {

	}
}
