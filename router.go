package mango

import (
	"net/http"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("\\{([\\w\\d]+)\\}")

//HandlerFunc use to handle incoming requests.
type HandlerFunc func(*Context)

//Route point to a URL path
type Route struct {
	Method      string
	Path        string
	Pathable    *regexp.Regexp
	Handler     HandlerFunc
	MiddlerPool []MiddlerFunc
	IsStatic    bool
}

//Router routes provider
type Router struct {
	StaticPool map[string][]*Route
	Pool       map[string][]*Route
}

//Push an route instance into pool
func (this *Router) Push(i *Route) {
	if i.IsStatic {
		this.StaticPool[i.Method] = append(this.StaticPool[i.Method], i)
	} else {
		this.Pool[i.Method] = append(this.Pool[i.Method], i)
	}
}

//Search matched route by given request instance
func (this *Router) Search(r *http.Request) (*Route, map[string]string) {
	var params map[string]string
	route := this.SearchStaticPool(r)
	if route == nil {
		route, params = this.SearchPool(r)
	}

	return route, params
}

//SearchStaticPool search route in static pool
func (this *Router) SearchStaticPool(r *http.Request) *Route {
	if batch, ok := this.StaticPool[r.Method]; ok {
		for _, route := range batch {
			if route.Path == r.RequestURI {
				return route
			}
		}
	}

	return nil
}

//SearchPool search route in custom pool
func (this *Router) SearchPool(r *http.Request) (*Route, map[string]string) {
	if batch, ok := this.Pool[r.Method]; ok {
		for _, route := range batch {
			if route.Pathable.MatchString(r.RequestURI) {
				params := map[string]string{}
				names := route.Pathable.SubexpNames()[1:]
				values := route.Pathable.FindStringSubmatch(r.RequestURI)[1:]

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

//CompilePath compile given path to regexp
//route path definition may with variables that defined
//as {uid}, it will compile to (?P<uid>[^/]+) and returns
//it as regexp.Regexp.
func (this *Router) CompilePath(path string) (*regexp.Regexp, bool) {
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

//Route register an new route to router.
func (this *Router) Route(methods []string, path string, fn HandlerFunc, middlers []MiddlerFunc) {
	path = strings.Trim(path, " /")
	path = "/" + path

	for _, m := range methods {
		p, s := this.CompilePath(path)
		route := &Route{
			Method:      m,
			Path:        path,
			Pathable:    p,
			Handler:     fn,
			MiddlerPool: middlers,
			IsStatic:    s,
		}

		this.Push(route)
	}
}

//Get register a GET route.
func (this *Router) Get(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"GET"}, path, fn, middlers)
}

//Post register a POST route.
func (this *Router) Post(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"POST"}, path, fn, middlers)
}

//Put register a PUT route.
func (this *Router) Put(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"PUT"}, path, fn, middlers)
}

//Delete register a DELETE route.
func (this *Router) Delete(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"DELETE"}, path, fn, middlers)
}

//Any register a route without request type limit.
func (this *Router) Any(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.Route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middlers)
}
