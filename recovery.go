package mango

import (
	"fmt"
	"net/http"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/common"
)

//Recovery record every panic.
func Recovery() common.MiddleFunc {
	return func(ctx common.Context) {
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

				ctx.Response().SetStatus(http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}
