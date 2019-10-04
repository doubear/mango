package contracts

import (
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

// Mango interface of mango micro framework.
type Mango interface {
	Any(string, Callable, ...ThenableFunc)
	Get(string, Callable, ...ThenableFunc)
	Post(string, Callable, ...ThenableFunc)
	Put(string, Callable, ...ThenableFunc)
	Delete(string, Callable, ...ThenableFunc)
	Group(string, func(Router), ...ThenableFunc)
	Use(ThenableFunc)
	SetDefaultRoute(Callable)
	SetCachable(Cachable)
	Start(string)
	StartTLS(string, string, string)
	StartAutoTLS(string, autocert.Cache, ...string)
	ServeHTTP(http.ResponseWriter, *http.Request)
	On(event string, fn func())
}
