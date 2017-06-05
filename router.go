package mango

import (
	"net/http"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("\\{([\\w\\d]+)\\}")

//HandlerFunc use to handle incoming requests.
type HandlerFunc func(*Context)

//route point to a URL path
type route struct {
	method      string
	path        string
	pathable    *regexp.Regexp
	handler     HandlerFunc
	middlerPool []MiddlerFunc
	isStatic    bool
}

//router routes provider
type router struct {
	StaticPool map[string][]*route
	Pool       map[string][]*route
}

//push an route instance into pool
func (this *router) push(i *route) {
	if i.isStatic {
		this.StaticPool[i.method] = append(this.StaticPool[i.method], i)
	} else {
		this.Pool[i.method] = append(this.Pool[i.method], i)
	}
}

//search matched route by given request instance
func (this *router) search(r *http.Request) (*route, map[string]string) {
	var params map[string]string
	route := this.searchStaticPool(r)
	if route == nil {
		route, params = this.searchPool(r)
	}

	return route, params
}

//searchStaticPool search route in static pool
func (this *router) searchStaticPool(r *http.Request) *route {
	if batch, ok := this.StaticPool[r.Method]; ok {
		for _, route := range batch {
			if route.path == r.RequestURI {
				return route
			}
		}
	}

	return nil
}

//searchPool search route in custom pool
func (this *router) searchPool(r *http.Request) (*route, map[string]string) {
	if batch, ok := this.Pool[r.Method]; ok {
		for _, route := range batch {
			if route.pathable.MatchString(r.RequestURI) {
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
func (this *router) compilePath(path string) (*regexp.Regexp, bool) {
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
func (this *router) route(methods []string, path string, fn HandlerFunc, middlers []MiddlerFunc) {
	path = strings.Trim(path, " /")
	path = "/" + path

	for _, m := range methods {
		p, s := this.compilePath(path)
		route := &route{
			method:      m,
			path:        path,
			pathable:    p,
			handler:     fn,
			middlerPool: middlers,
			isStatic:    s,
		}

		this.push(route)
	}
}

//Get register a GET route.
func (this *router) Get(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.route([]string{"GET"}, path, fn, middlers)
}

//Post register a POST route.
func (this *router) Post(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.route([]string{"POST"}, path, fn, middlers)
}

//Put register a PUT route.
func (this *router) Put(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.route([]string{"PUT"}, path, fn, middlers)
}

//Delete register a DELETE route.
func (this *router) Delete(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.route([]string{"DELETE"}, path, fn, middlers)
}

//Any register a route without request type limit.
func (this *router) Any(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middlers)
}
