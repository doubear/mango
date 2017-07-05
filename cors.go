package mango

//CorsOption configure the cors middleware.
type CorsOption struct {
	Origin  string
	Methods string
	Headers string
}

//Cors additional CORS middleware
func Cors(opt CorsOption) MiddleFunc {
	return func(ctx *Context) {
		// if ctx.R.Method == "OPTIONS" {
		ctx.W.Header().Add("Access-Control-Allow-Origin", opt.Origin)
		ctx.W.Header().Add("Access-Control-Allow-Methods", opt.Methods)
		ctx.W.Header().Add("Access-Control-Allow-Headers", opt.Headers)
		// ctx.W.SetStatus(http.StatusOK)
		// } else {
		ctx.Next()
		// }
	}
}
