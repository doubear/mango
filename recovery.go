package mango

import (
	"fmt"
	"net/http"

	"github.com/go-mango/logy"
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

				logy.W(err.Error())

				ctx.W.SetStatus(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}
