package concretes

import (
	"regexp"

	"github.com/go-mango/mango/contracts"
)

var re = regexp.MustCompile("\\{([\\w\\d]+)\\}")

type route struct {
	method    string
	path      string
	pathable  *regexp.Regexp
	callable  contracts.Callable
	thenStack []contracts.ThenableFunc
	isStatic  bool
}

// NewRoute returns route instance.
func NewRoute(method string, path string, callable contracts.Callable, stack ...contracts.ThenableFunc) contracts.Route {
	pathable, isStatic := compilePath(path)

	return &route{
		method,
		path,
		pathable,
		callable,
		stack,
		isStatic,
	}
}

//compilePath compile given path to regexp
//route path definition may with variables that defined
//as {uid}, it will compile to (?P<uid>[^/]+) and returns
//it as regexp.Regexp.
func compilePath(path string) (*regexp.Regexp, bool) {
	if path == "" {
		path = "/"
	}

	pathen := re.ReplaceAllString(path, "(?P<$1>[^/]+)")

	is := pathen == path

	if is == false {
		pathen = "^" + pathen + "$"
	}

	return regexp.MustCompile(pathen), is
}

func (route *route) SetIsStatic(state bool) {
	route.isStatic = state
}

func (route *route) SetThenStack(stack ...contracts.ThenableFunc) {
	route.thenStack = append(route.thenStack, stack...)
}

func (route *route) Method() string {
	return route.method
}

func (route *route) Path() string {
	return route.path
}

func (route *route) SetPath(path string) {
	route.path = path
}

func (route *route) Pathable() *regexp.Regexp {
	return route.pathable
}

func (route *route) Callable() contracts.Callable {
	return route.callable
}

func (route *route) ThenStack() []contracts.ThenableFunc {
	return route.thenStack
}

func (route *route) IsStatic() bool {
	return route.isStatic
}
