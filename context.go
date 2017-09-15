package mango

import (
	"strings"

	"github.com/go-mango/mango/common"
)

type context struct {
	R       common.Request
	W       common.Response
	C       Cacher
	middles []common.MiddleFunc
	dict    map[string]interface{}
}

func newContext(r common.Request, w common.Response, c Cacher) common.Context {
	return &context{
		r,
		w,
		c,
		make([]common.MiddleFunc, 0),
		make(map[string]interface{}),
	}
}

//Request returns wrapped http.Request
func (c *context) Request() common.Request {
	return c.R
}

//Response returns wrapped http.ResponseWriter
func (c *context) Response() common.Response {
	return c.W
}

//Next executes the next middleware func
func (c *context) Next() {
	if len(c.middles) > 0 {
		m := c.middles[0]
		c.middles = c.middles[1:]
		m(c)
	}
}

//Get retrieves an temporary variable.
func (c *context) Get(name string) interface{} {
	if v, ok := c.dict[name]; ok {
		return v
	}

	return nil
}

//Set stores a key-value pair.
func (c *context) Set(name string, value interface{}) {
	c.dict[name] = value
}

//URL generates URL with given params.
func (c *context) URL(u string, p map[string]string) string {
	for _, k := range p {
		u = strings.Replace(u, "{"+k+"}", p[k], -1)
	}

	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

	return u
}
