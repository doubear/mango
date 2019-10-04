package concretes

import (
	"strings"

	"github.com/go-mango/mango/contracts"
)

type context struct {
	request  contracts.Request
	response contracts.Response
	cache    contracts.Cachable
	stack    []contracts.ThenableFunc
	auth     contracts.Authenable
	session  contracts.Session
}

// NewContext create new Context instance
func NewContext(
	request contracts.Request,
	response contracts.Response,
	cache contracts.Cachable,
	stack []contracts.ThenableFunc,
	route contracts.Route,
) contracts.ThenableContext {
	return &context{
		request,
		response,
		cache,
		append(stack, handleResponse(route.Callable())),
		newAuth(),
		NewSession(),
	}
}

// Request returns wrapped http.Request
func (c *context) Request() contracts.Request {
	return c.request
}

// Response returns wrapped http.ResponseWriter
func (c *context) Response() contracts.Response {
	return c.response
}

// Next executes the next middleware func
func (c *context) Next() {
	if len(c.stack) > 0 {
		m := c.stack[0]
		c.stack = c.stack[1:]
		m(c)
	}
}

// URL generates URL with given params.
func (c *context) URL(u string, p map[string]string) string {
	for _, k := range p {
		u = strings.Replace(u, "{"+k+"}", p[k], -1)
	}

	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

	return u
}

// Auth returns auth provider of incoming request.
func (c *context) Auth() contracts.Authenable {
	return c.auth
}

// Cache returns cache storage instance of incoming request.
func (c *context) Cache() contracts.Cachable {
	return c.cache
}

// Session returns session storage instace of incoming request.
func (c *context) Session() contracts.Session {
	return c.session
}
