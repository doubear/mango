package mango

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/contracts"
)

//Recovery record every panic.
func Recovery() contracts.MiddleFunc {
	return func(ctx contracts.Context) {
		defer func() {
			if v := recover(); v != nil {
				var err error
				switch v := v.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", v)
				}

				logy.Std().Warn(err.Error())

				debug.PrintStack()

				ctx.Response().SetStatus(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}
