package mango

import "net/http"

//Cors addtional CORS middleware
func Cors() MiddleFunc {
	return func(ctx *Context) {
		if ctx.R.Method == "OPTIONS" {
			ctx.W.Header().Add("Access-Control-Allow-Origin", "*")
			ctx.W.Header().Add("Access-Control-Allow-Methods", "*")
			ctx.W.Header().Add("Access-Control-Allow-Headers", "*")
			ctx.W.Header().Add("Access-Control-Max-Age", "86400")
			ctx.W.SetStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
