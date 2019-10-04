package middlewares

import (
	"net/http"

	"github.com/go-mango/mango/contracts"
)

//RedirectOption configures HTTPS Redirector.
type RedirectOption struct {
	MustHTTPS bool
	MustHOST  string
}

//Redirect recirects  requests.
func Redirect(opt RedirectOption) contracts.ThenableFunc {
	return func(ctx contracts.ThenableContext) {
		if opt.MustHOST != "" && ctx.Request().Host() != opt.MustHOST {
			to := *ctx.Request().URL()
			to.Scheme = "http"
			to.Host = opt.MustHOST

			ctx.Response().Redirect(http.StatusPermanentRedirect, to.String())
			return
		}

		if opt.MustHTTPS && false == ctx.Request().IsTLS() {
			to := *ctx.Request().URL()
			to.Scheme = "https"
			to.Host = ctx.Request().Host()

			ctx.Response().Redirect(http.StatusPermanentRedirect, to.String())
			return
		}

		ctx.Next()
	}
}
