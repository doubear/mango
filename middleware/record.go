package middleware

import (
	"time"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango"
)

//Record log incoming requests to console.
func Record() mango.MiddleFunc {
	return func(ctx mango.Context) {
		st := time.Now()
		ctx.Next()
		dur := time.Since(st).String()

		logy.I(
			"%s %s %s %d %dB %s",
			ctx.Request().IP(),
			ctx.Request().Method(),
			ctx.Request().URI(),
			ctx.Response().Status(),
			ctx.Response().Size(),
			dur,
		)
	}
}
