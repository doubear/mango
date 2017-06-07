package mango

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"os"

	"github.com/go-mango/mango/logger"
)

const proxymark = "X-Forwarded-For"

//MiddleFunc used as a middleware
type MiddleFunc func(*Context)

//Context income request context
type Context struct {
	R       *http.Request
	W       *response
	params  map[string]string
	middles []MiddleFunc
	Logger  *logger.Logger
}

//Next executes the next middleware func
func (ctx *Context) Next() {
	if len(ctx.middles) > 0 {
		m := ctx.middles[0]
		ctx.middles = ctx.middles[1:]
		m(ctx)
	}
}

//ClientIP returns connected client's IP address.
func (ctx *Context) ClientIP() string {
	ip := ctx.R.RemoteAddr

	if ctx.R.Header.Get(proxymark) != "" {
		//using proxy server
		proxy := strings.Split(ctx.R.Header.Get(proxymark), ",")[0]
		proxy = strings.TrimSpace(proxy)
		proxyIP := net.ParseIP(proxy)
		if false == proxyIP.IsGlobalUnicast() {
			ip = proxyIP.String()
		}
	}

	ip = strings.Split(ip, ":")[0] //to fixed r.RemoteAddr format.

	return ip
}

//SaveFile receives file from MULTI-PART FORM.
func (ctx *Context) SaveFile(field string, saveTo io.Writer) (string, bool) {
	f, h, err := ctx.R.FormFile(field)
	if err != nil {
		return "", false
	}

	_, err = io.Copy(saveTo, f)
	if err != nil {
		return "", false
	}

	return h.Filename, true
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

//JSON auto-encode given value to json then write it to response.
func (ctx *Context) JSON(code int, v interface{}) {
	ctx.W.SetStatus(code)
	e := json.NewEncoder(ctx.W)
	e.Encode(v)
}

//Blob writes given bytes to response.
func (ctx *Context) Blob(code int, b []byte) {
	ctx.W.SetStatus(code)
	ctx.W.Write(b)
}

//Text writes given string to response.
func (ctx *Context) Text(code int, s string) {
	ctx.W.SetStatus(code)
	ctx.W.WriteString(s)
}

//File serves file from given path.
func (ctx *Context) File(path string) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.W.SetStatus(http.StatusNotFound)
			ctx.W.Clear()
			return
		}

		ctx.W.SetStatus(http.StatusInternalServerError)
		ctx.W.Clear()
		return
	}

	defer file.Close()

	io.Copy(ctx.W, file)
}
