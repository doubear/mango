package mango

import (
	"strings"

	"github.com/go-mango/mango/contracts"
)

//GroupFunc use to register route group.
type GroupFunc func(*GroupRouter)

//GroupRouter register routes in group.
type GroupRouter struct {
	paths   []string
	middles []contracts.MiddleFunc
	router  *router
}

func (grt *GroupRouter) path(s string) string {
	paths := grt.paths
	s = strings.Trim(s, " /")

	if s != "" {
		paths = append(grt.paths, s)
	}

	return "/" + strings.Join(paths, "/")
}

func (grt *GroupRouter) middler(ms []contracts.MiddleFunc) []contracts.MiddleFunc {
	ms = append(grt.middles, ms...)

	return ms
}

//Use register group router middleware.
func (grt *GroupRouter) Use(m contracts.MiddleFunc) {
	grt.middles = append(grt.middles, m)
}

//Route appends group route to base router.
func (grt *GroupRouter) Route(methods []string, path string, fn contracts.HandlerFunc, ms []contracts.MiddleFunc) {
	path = grt.path(path)
	ms = grt.middler(ms)

	grt.router.route(methods, path, fn, ms)
}

//Get register a GET route.
func (grt *GroupRouter) Get(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	grt.Route([]string{"GET"}, path, fn, middles)
}

//Post register a POST route.
func (grt *GroupRouter) Post(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	grt.Route([]string{"POST"}, path, fn, middles)
}

//Put register a PUT route.
func (grt *GroupRouter) Put(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	grt.Route([]string{"PUT"}, path, fn, middles)
}

//Delete register a DELETE route.
func (grt *GroupRouter) Delete(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	grt.Route([]string{"DELETE"}, path, fn, middles)
}

//Any register a route without request type limit.
func (grt *GroupRouter) Any(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	grt.Route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middles)
}

//Group create an new subgroup.
func (grt *GroupRouter) Group(path string, fn GroupFunc, middles ...contracts.MiddleFunc) {
	paths := grt.paths

	path = strings.Trim(path, " /")
	if path != "" {
		paths = append(paths, path)
	}

	middles = append(grt.middles, middles...)

	fn(&GroupRouter{
		paths,
		middles,
		grt.router,
	})
}
