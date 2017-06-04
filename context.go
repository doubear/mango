package mango

import "net/http"
import "strings"
import "net"

//MiddlerFunc used as a middleware
type MiddlerFunc func(*Context)

//Context income request context
type Context struct {
	*http.Request
	W        *response
	params   map[string]string
	middlers []MiddlerFunc
	IP       string
	ViaProxy bool //shows the request is sent via proxy server.
	Logger   *Logger
}

//Param get value of path param
func (this *Context) Param(k, d string) string {
	if v, ok := this.params[k]; ok {
		return v
	}

	return d
}

//Next executes the next middleware func
func (this *Context) Next() {
	m := this.middlers[0]
	this.middlers = this.middlers[1:]
	m(this)
}

//NewContext create new Context instance
func newContext(r *http.Request, w http.ResponseWriter, ps map[string]string, ms []MiddlerFunc) *Context {
	ip, isProxy := getRealClientIP(r)

	return &Context{
		r,
		newResponse(w),
		ps,
		ms,
		ip,
		isProxy,
		NewLogger(),
	}
}

const proxymark = "X-Forwarded-For"

func getRealClientIP(r *http.Request) (string, bool) {
	ip := r.RemoteAddr
	is := false

	if r.Header.Get(proxymark) != "" {
		//using proxy server
		proxy := strings.Split(r.Header.Get(proxymark), ",")[0]
		proxy = strings.TrimSpace(proxy)
		proxyIP := net.ParseIP(proxy)
		if false == proxyIP.IsGlobalUnicast() {
			ip = proxyIP.String()
			is = true
		}
	}

	ip = strings.Split(ip, ":")[0] //to fixed r.RemoteAddr format.

	return ip, is
}
