package mango

import "net/http"
import "strings"
import "net"
import "io"
import "encoding/json"

const proxymark = "X-Forwarded-For"

//MiddlerFunc used as a middleware
type MiddlerFunc func(*Context)

//Context income request context
type Context struct {
	R        *http.Request
	W        *response
	params   map[string]string
	middlers []MiddlerFunc
	Logger   *Logger
}

//Next executes the next middleware func
func (this *Context) Next() {
	if len(this.middlers) > 0 {
		m := this.middlers[0]
		this.middlers = this.middlers[1:]
		m(this)
	}
}

//ClientIP returns connected client's IP address.
func (this *Context) ClientIP() string {
	ip := this.R.RemoteAddr

	if this.R.Header.Get(proxymark) != "" {
		//using proxy server
		proxy := strings.Split(this.R.Header.Get(proxymark), ",")[0]
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
func (this *Context) File(field string, saveTo io.Writer) (string, bool) {
	f, h, err := this.R.FormFile(field)
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
func (this *Context) Form(field string) string {
	return this.R.PostFormValue(field)
}

//Query retrieves value from GET params.
func (this *Context) Query(field string) string {
	return this.R.URL.Query().Get(field)
}

//Param retrieves value from PATH params.
func (this *Context) Param(k, d string) string {
	if v, ok := this.params[k]; ok {
		return v
	}

	return d
}

//Input retrieves value with given field name from both Form and Query.
func (this *Context) Input(field string) string {
	if v := this.Form(field); v != "" {
		return v
	}

	return this.Query(field)
}

//JSON auto-encode given value to json then write it to response.
func (this *Context) JSON(code int, v interface{}) {
	this.W.SetStatus(code)
	e := json.NewEncoder(this.W)
	e.Encode(v)
}

//Blob writes given bytes to response.
func (this *Context) Blob(code int, b []byte) {
	this.W.SetStatus(code)
	this.W.Write(b)
}

//Text writes given string to response.
func (this *Context) Text(code int, s string) {
	this.W.SetStatus(code)
	this.W.WriteString(s)
}

//NewContext create new Context instance
func newContext(r *http.Request, w http.ResponseWriter, ps map[string]string, ms []MiddlerFunc) *Context {
	return &Context{
		r,
		newResponse(w),
		ps,
		ms,
		NewLogger(),
	}
}
