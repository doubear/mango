package mango

import (
	"fmt"
	"net/http"
)

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
