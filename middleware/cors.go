package middleware

import (
	"net/http"

	"github.com/go-mango/mango/common"
)

//CorsOption configure the cors middleware.
type CorsOption struct {
	Origin  string
	Methods string
	Headers string
}

//Cors additional CORS middleware
func Cors(opt CorsOption) common.MiddleFunc {
	return func(ctx common.Context) {
		ctx.Response().Header().Add("Access-Control-Allow-Origin", opt.Origin)
		ctx.Response().Header().Add("Access-Control-Allow-Methods", opt.Methods)
		ctx.Response().Header().Add("Access-Control-Allow-Headers", opt.Headers)

		if ctx.Request().Method() == "OPTIONS" {
			ctx.Response().SetStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
