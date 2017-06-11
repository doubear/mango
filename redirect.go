package mango

import (
	"net/http"
)

//RedirectOption configures HTTPS Redirector.
type RedirectOption struct {
	MustHTTPS bool
	MustHOST  string
}

//Redirect recirects  requests.
func Redirect(opt RedirectOption) MiddleFunc {
	return func(ctx *Context) {
		if opt.MustHOST != "" && ctx.R.Host != opt.MustHOST {
			to := *ctx.R.URL
			to.Scheme = "http"
			to.Host = opt.MustHOST

			ctx.W.Redirect(http.StatusPermanentRedirect, to.String())
			return
		}

		if opt.MustHTTPS && false == ctx.IsTLS() {
			to := *ctx.R.URL
			to.Scheme = "https"
			to.Host = ctx.R.Host

			ctx.W.Redirect(http.StatusPermanentRedirect, to.String())
			return
		}

		ctx.Next()
	}
}
