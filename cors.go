package mango

import "net/http"

//CorsOption configure the cors middleware.
type CorsOption struct {
	Origin  string
	Methods string
	Headers string
}

//Cors additional CORS middleware
func Cors(opt CorsOption) MiddleFunc {
	return func(ctx *Context) {
		ctx.W.Header().Add("Access-Control-Allow-Origin", opt.Origin)
		ctx.W.Header().Add("Access-Control-Allow-Methods", opt.Methods)
		ctx.W.Header().Add("Access-Control-Allow-Headers", opt.Headers)

		if ctx.R.Method == "OPTIONS" {
			ctx.W.SetStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
