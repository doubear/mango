package mango

import (
	"net/http"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("\\{([\\w\\d]+)\\}")

//HandlerFunc use to handle incoming requests.
type HandlerFunc func(Context) (int, interface{})

//route point to a URL path
type route struct {
	method     string
	path       string
	pathable   *regexp.Regexp
	handler    HandlerFunc
	middlePool []MiddleFunc
	isStatic   bool
}

//router routes provider
type router struct {
	StaticPool map[string][]*route
	Pool       map[string][]*route
}

//push an route instance into pool
func (rt *router) push(i *route) {
	if i.isStatic {
		rt.StaticPool[i.method] = append(rt.StaticPool[i.method], i)
	} else {
		rt.Pool[i.method] = append(rt.Pool[i.method], i)
	}
}

//search matched route by given request instance
func (rt *router) search(r *http.Request) (*route, map[string]string) {
	var params map[string]string
	route := rt.searchStaticPool(r)
	if route == nil {
		route, params = rt.searchPool(r)
	}

	return route, params
}

//searchStaticPool search route in static pool
func (rt *router) searchStaticPool(r *http.Request) *route {
	if batch, ok := rt.StaticPool[r.Method]; ok {
		for _, route := range batch {
			if route.path == r.URL.Path {
				return route
			}
		}
	}

	return nil
}

//searchPool search route in custom pool
func (rt *router) searchPool(r *http.Request) (*route, map[string]string) {
	if batch, ok := rt.Pool[r.Method]; ok {
		for _, route := range batch {
			if route.pathable.MatchString(r.URL.Path) {
				params := map[string]string{}
				names := route.pathable.SubexpNames()[1:]
				values := route.pathable.FindStringSubmatch(r.RequestURI)[1:]

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

//compilePath compile given path to regexp
//route path definition may with variables that defined
//as {uid}, it will compile to (?P<uid>[^/]+) and returns
//it as regexp.Regexp.
func (rt *router) compilePath(path string) (*regexp.Regexp, bool) {
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

//route register an new route to router.
func (rt *router) route(methods []string, path string, fn HandlerFunc, middlers []MiddleFunc) {
	path = strings.Trim(path, " /")
	path = "/" + path

	for _, m := range methods {
		p, s := rt.compilePath(path)
		route := &route{
			method:     m,
			path:       path,
			pathable:   p,
			handler:    fn,
			middlePool: middlers,
			isStatic:   s,
		}

		rt.push(route)
	}
}

//Get register a GET route.
func (rt *router) Get(path string, fn HandlerFunc, middlers ...MiddleFunc) {
	rt.route([]string{"GET"}, path, fn, middlers)
}

//Post register a POST route.
func (rt *router) Post(path string, fn HandlerFunc, middlers ...MiddleFunc) {
	rt.route([]string{"POST"}, path, fn, middlers)
}

//Put register a PUT route.
func (rt *router) Put(path string, fn HandlerFunc, middlers ...MiddleFunc) {
	rt.route([]string{"PUT"}, path, fn, middlers)
}

//Delete register a DELETE route.
func (rt *router) Delete(path string, fn HandlerFunc, middlers ...MiddleFunc) {
	rt.route([]string{"DELETE"}, path, fn, middlers)
}

//Any register a route without request type limit.
func (rt *router) Any(path string, fn HandlerFunc, middlers ...MiddleFunc) {
	rt.route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middlers)
}
