package common

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"io/ioutil"
	"github.com/go-mango/mango/contracts"

	"encoding/json"
)

type request struct {
	reader *http.Request
	params map[string]string
}

//NewRequest create new request instance.
func NewRequest(r *http.Request, ps map[string]string) contracts.Request {
	return &request{
		r,
		ps,
	}
}

//Parent returns original http.Request
func (c *request) Parent() *http.Request {
	return c.reader
}

//IP returns connected client's IP address.
func (c *request) IP() string {
	ip := c.reader.RemoteAddr

	if c.reader.Header.Get("X-Forwarded-For") != "" {
		//using proxy server
		proxy := strings.Split(c.reader.Header.Get("X-Forwarded-For"), ",")[0]
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
func (c *request) File(k string) (contracts.UploadedFile, error) {
	f, h, err := c.reader.FormFile(k)
	if err != nil {
		return nil, err
	}

	return NewUploadedFile(h, f), nil
}

//Form retrieves value from POST form.
func (c *request) Form(k string) string {
	return c.reader.PostFormValue(k)
}

//Query retrieves value from GET params.
func (c *request) Query(k string) string {
	return c.reader.URL.Query().Get(k)
}

//Param retrieves value from PATH params.
func (c *request) Param(k string) string {
	if v, ok := c.params[k]; ok {
		return v
	}

	return ""
}

//Input retrieves value with given k name from both Form and Query.
func (c *request) Input(k string) string {
	if v := c.Form(k); v != "" {
		return v
	}

	return c.Query(k)
}

//JSON parse request body as JSON.
func (c *request) JSON(v interface{}) error {
	data, err := ioutil.ReadAll(c.reader.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

//IsTLS detects the request is over HTTPS or not.
func (c *request) IsTLS() bool {
	return c.reader.TLS != nil
}

//Header returns original http.Header.
func (c *request) Header() http.Header {
	return c.reader.Header
}

//Method returns request method.
func (c *request) Method() string {
	return c.reader.Method
}

//URI returns http.Request.RequestURI
func (c *request) URI() string {
	return c.reader.RequestURI
}

func (c *request) URL() *url.URL {
	return c.reader.URL
}

func (c *request) Host() string {
	return c.reader.Host
}
