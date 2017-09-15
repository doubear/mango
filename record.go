package mango

import (
	"time"

	"github.com/go-mango/logy"
)

//Record log incoming requests to console.
func Record() MiddleFunc {
	return func(ctx *Context) {
		st := time.Now()
		ctx.Next()
		dur := time.Since(st).String()

		logy.I(
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
