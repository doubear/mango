package mango

import (
	"time"
)

//Record log incoming requests to console.
func Record() MiddleFunc {
	return func(ctx *Context) {
		st := time.Now()
		ctx.Next()
		dur := NumericTimeSmartFormat(time.Since(st).Nanoseconds())

		ctx.Logger.Infof(
			"%s %s %s %d %dB %s",
			ctx.ClientIP(),
			ctx.R.Method,
			ctx.R.RequestURI,
			ctx.W.Status(),
			ctx.W.Size(),
			dur,
		)
	}
}
