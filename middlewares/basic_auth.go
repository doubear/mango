package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/go-mango/mango/contracts"
)

//BasicAuth provides basic-auth middleware.
//
/*
	credentials = map[string]string{
		"username": "password",
	}

	m.Use(mango.BasicAuth(credentials))
*/
func BasicAuth(credentials map[string]string) contracts.ThenableFunc {
	return func(ctx contracts.ThenableContext) {
		if token := ctx.Request().Header().Get("Authorization"); token != "" {
			if strings.HasPrefix(token, "Basic ") {
				token = token[6:]

				raw, err := base64.StdEncoding.DecodeString(token)
				if err != nil {
					ctx.Response().SetStatus(http.StatusInternalServerError)
					return
				}

				credential := strings.SplitN(string(raw), ":", 1)

				if pwd, ok := credentials[credential[0]]; ok {
					if pwd == credential[1] {
						ctx.Next()
						return
					}
				}
			}
		}

		ctx.Response().SetStatus(http.StatusUnauthorized)
		ctx.Response().Header().Set("WWW-Authenticate", "Basic realm=\"Restricted\"")
	}
}
