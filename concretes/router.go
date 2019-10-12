package concretes

import (
	"strings"

	"github.com/go-mango/mango/contracts"
)

type router struct {
	prefixes     []string
	stack        []contracts.ThenableFunc
	staticPool   map[string][]contracts.Route
	pool         map[string][]contracts.Route
	defaultRoute contracts.Route
}

var defaultRoute = NewRoute("*", "/", func(ctx contracts.Context) (int, interface{}) {
	return 404, nil
})

// NewRouter create new router.
func NewRouter() contracts.Router {
	return &router{
		[]string{""},
		[]contracts.ThenableFunc{},
		map[string][]contracts.Route{},
		map[string][]contracts.Route{},
		defaultRoute,
	}
}

func (router *router) SetDefaultRoute(callable contracts.Callable) {
	router.defaultRoute = NewRoute("*", "/", callable)
}

func (router *router) push(route contracts.Route) {
	if route.IsStatic() {
		router.staticPool[route.Method()] = append(router.staticPool[route.Method()], route)
	} else {
		router.pool[route.Method()] = append(router.pool[route.Method()], route)
	}
}

func (router *router) ToMatch(r contracts.Request) (contracts.Route, map[string]string) {
	var params map[string]string
	route := router.searchStaticPool(r)
	if route == nil {
		route, params = router.searchPool(r)
	}

	if route == nil {
		route = router.defaultRoute
	}

	return route, params
}

func (router *router) searchStaticPool(r contracts.Request) contracts.Route {
	if batch, ok := router.staticPool[r.Method()]; ok {
		for _, route := range batch {
			if route.Path() == r.URL().Path {
				return route
			}
		}
	}

	return nil
}

//searchPool search route in custom pool
func (router *router) searchPool(r contracts.Request) (contracts.Route, map[string]string) {
	if batch, ok := router.pool[r.Method()]; ok {
		for _, route := range batch {
			if route.Pathable().MatchString(r.URL().Path) {
				params := map[string]string{}
				names := route.Pathable().SubexpNames()[1:]
				values := route.Pathable().FindStringSubmatch(r.URL().Path)[1:]

				if len(names) != len(values) {
					continue
				}

				for i, name := range names {
					params[name] = values[i]
				}

				return route, params
			}
		}
	}

	return nil, nil
}

func (router *router) Use(next ...contracts.ThenableFunc) {
	router.stack = append(router.stack, next...)
}

func (router *router) Prefixes() []string {
	return router.prefixes
}

func (router *router) ThenableStack() []contracts.ThenableFunc {
	return router.stack
}

func (router *router) SetThenableStack(stack ...contracts.ThenableFunc) {
	router.stack = stack
}

// Group performs batch routes registration with same URI prefix.
func (router *router) Group(prefix string, entry func(contracts.Router), stack ...contracts.ThenableFunc) {
	router.pushScope(prefix)
	savedStack := router.ThenableStack()
	router.Use(router.stack...)

	entry(router)

	router.SetThenableStack(savedStack...)
	router.popScope()
}

func (router *router) pushScope(scope string) {
	scope = strings.Trim(scope, " /")
	router.prefixes = append(router.prefixes, scope)
}

func (router *router) popScope() {
	router.prefixes = router.prefixes[:len(router.prefixes)-1]
}

func (router *router) newScopedRoute(
	method string,
	path string,
	resolver contracts.Callable,
	stack ...contracts.ThenableFunc,
) {
	router.pushScope(path)
	path = strings.Join(router.prefixes, "/")
	stack = append(router.stack, stack...)
	router.push(NewRoute(method, path, resolver, stack...))
	router.popScope()
}

// Any register resolver function for route prefixed with "prefix".
func (router *router) Any(path string, resolver contracts.Callable, stack ...contracts.ThenableFunc) {
	router.Get(path, resolver, stack...)
	router.Post(path, resolver, stack...)
	router.Put(path, resolver, stack...)
	router.Delete(path, resolver, stack...)
}

// Get register resolver function called by GET requests.
func (router *router) Get(path string, resolver contracts.Callable, stack ...contracts.ThenableFunc) {
	router.newScopedRoute("GET", path, resolver, stack...)
}

// Post register resolver function called by POST requests.
func (router *router) Post(path string, resolver contracts.Callable, stack ...contracts.ThenableFunc) {
	router.newScopedRoute("POST", path, resolver, stack...)
}

// Put register resolver function called by PUT requests.
func (router *router) Put(path string, resolver contracts.Callable, stack ...contracts.ThenableFunc) {
	router.newScopedRoute("PUT", path, resolver, stack...)
}

// Delete register resolver function called by DELETE requests.
func (router *router) Delete(path string, resolver contracts.Callable, stack ...contracts.ThenableFunc) {
	router.newScopedRoute("DELETE", path, resolver, stack...)
}
