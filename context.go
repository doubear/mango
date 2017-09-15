package mango

import (
	"io"
	"net"
	"net/http"
	"strings"

	"io/ioutil"

	"encoding/json"

	"mime/multipart"
)

//MiddleFunc used as a middleware
type MiddleFunc func(*Context)

//Context income request context
type Context struct {
	R       *http.Request
	W       *response
	C       Cacher
	params  map[string]string
	middles []MiddleFunc
	dict    map[string]interface{}
}

//Next executes the next middleware func
func (ctx *Context) Next() {
	if len(ctx.middles) > 0 {
		m := ctx.middles[0]
		ctx.middles = ctx.middles[1:]
		m(ctx)
	}
}

//Get retrieves an temporary variable.
func (ctx *Context) Get(name string) interface{} {
	if v, ok := ctx.dict[name]; ok {
		return v
	}

	return nil
}

//Set stores a key-value pair.
func (ctx *Context) Set(name string, value interface{}) {
	ctx.dict[name] = value
}

//ClientIP returns connected client's IP address.
func (ctx *Context) ClientIP() string {
	ip := ctx.R.RemoteAddr

	if ctx.R.Header.Get("X-Forwarded-For") != "" {
		//using proxy server
		proxy := strings.Split(ctx.R.Header.Get("X-Forwarded-For"), ",")[0]
		proxy = strings.TrimSpace(proxy)
		proxyIP := net.ParseIP(proxy)
		if false == proxyIP.IsGlobalUnicast() {
			ip = proxyIP.String()
		}
	}

	ip = strings.Split(ip, ":")[0] //to fixed r.RemoteAddr format.

	return ip
}

//File receives file from MULTI-PART FORM.
func (ctx *Context) File(field string, saveTo io.Writer) (*multipart.FileHeader, bool) {
	f, h, err := ctx.R.FormFile(field)
	if err != nil {
		return nil, false
	}

	_, err = io.Copy(saveTo, f)
	if err != nil {
		return nil, false
	}

	return h, true
}

//Form retrieves value from POST form.
func (ctx *Context) Form(field string) string {
	return ctx.R.PostFormValue(field)
}

//Query retrieves value from GET params.
func (ctx *Context) Query(field string) string {
	return ctx.R.URL.Query().Get(field)
}

//Param retrieves value from PATH params.
func (ctx *Context) Param(k, d string) string {
	if v, ok := ctx.params[k]; ok {
		return v
	}

	return d
}

//Input retrieves value with given field name from both Form and Query.
func (ctx *Context) Input(field string) string {
	if v := ctx.Form(field); v != "" {
		return v
	}

	return ctx.Query(field)
}

//JSON parse request body as JSON.
func (ctx *Context) JSON(v interface{}) error {
	data, err := ioutil.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

//IsTLS returns request is over HTTPS or not.
func (ctx *Context) IsTLS() bool {
	return ctx.R.TLS != nil
}

//URL generates URL with given params.
func (ctx *Context) URL(u string, p map[string]string) string {
	for _, k := range p {
		u = strings.Replace(u, "{"+k+"}", p[k], -1)
	}

	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

	return u
}
