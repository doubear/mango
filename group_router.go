package mango

import (
	"strings"
)

//GroupFunc use to register route group.
type GroupFunc func(*GroupRouter)

//GroupRouter register routes in group.
type GroupRouter struct {
	paths    []string
	middlers []MiddlerFunc
	router   *Router
}

func (this *GroupRouter) path(s string) string {
	paths := this.paths
	s = strings.Trim(s, " /")

	if s != "" {
		paths = append(this.paths, s)
	}

	return "/" + strings.Join(paths, "/")
}

func (this *GroupRouter) middler(ms []MiddlerFunc) []MiddlerFunc {
	ms = append(this.middlers, ms...)

	return ms
}

//Use register group router middleware.
func (this *GroupRouter) Use(m MiddlerFunc) {
	this.middlers = append(this.middlers, m)
}

func (this *GroupRouter) Route(methods []string, path string, fn HandlerFunc, ms []MiddlerFunc) {
	path = this.path(path)
	ms = this.middler(ms)

	this.router.Route(methods, path, fn, ms)
}

//Get register a GET route.
func (this *GroupRouter) Get(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"GET"}, path, fn, middlers)
}

//Post register a POST route.
func (this *GroupRouter) Post(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"POST"}, path, fn, middlers)
}

//Put register a PUT route.
func (this *GroupRouter) Put(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"PUT"}, path, fn, middlers)
}

//Delete register a DELETE route.
func (this *GroupRouter) Delete(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"DELETE"}, path, fn, middlers)
}

//Any register a route without request type limit.
func (this *GroupRouter) Any(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middlers)
}

//Group create an new subgroup.
func (this *GroupRouter) Group(path string, fn GroupFunc, middlers ...MiddlerFunc) {
	paths := this.paths

	path = strings.Trim(path, " /")
	if path != "" {
		paths = append(paths, path)
	}

	middlers = append(this.middlers, middlers...)

	fn(&GroupRouter{
		paths,
		middlers,
		this.router,
	})
}
