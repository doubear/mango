package middlewares

import (
	"time"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/contracts"
)

//Record log incoming requests to console.
func Record() contracts.ThenableFunc {
	return func(ctx contracts.ThenableContext) {
		st := time.Now()
		ctx.Next()
		dur := time.Since(st).String()

		logy.Std().Infof(
			"%d %s\t%dB\t%s\t%s\t%s",
			ctx.Response().Status(),
			dur,
			ctx.Response().Size(),
			ctx.Request().IP(),
			ctx.Request().Method(),
			ctx.Request().URI(),
		)
	}
}
