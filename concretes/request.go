package concretes

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
	parent *http.Request
	args   map[string]string
}

// NewRequest create new request instance.
func NewRequest(parent *http.Request) contracts.Request {
	return &request{
		parent,
		map[string]string{},
	}
}

// SetArgs set arguments for incoming requests that parsed from itself.
func (request *request) SetArgs(args map[string]string) {
	request.args = args
}

// Args returns all argument of incoming requests.
func (request *request) Args() map[string]string {
	return request.args
}

// Parent returns original http.Request
func (request *request) Parent() *http.Request {
	return request.parent
}

// IP returns connected client's IP address.
func (request *request) IP() string {
	ip := request.parent.RemoteAddr

	if request.parent.Header.Get("X-Forwarded-For") != "" {
		//using proxy server
		proxy := strings.Split(request.parent.Header.Get("X-Forwarded-For"), ",")[0]
		proxy = strings.TrimSpace(proxy)
		proxyIP := net.ParseIP(proxy)
		if false == proxyIP.IsGlobalUnicast() {
			ip = proxyIP.String()
		}
	}

	ip = strings.Split(ip, ":")[0] //to fixed r.RemoteAddr format.

	return ip
}

// File receives file from MULTI-PART FORM.
func (request *request) File(k string) (contracts.UploadedFile, error) {
	f, h, err := request.parent.FormFile(k)
	if err != nil {
		return nil, err
	}

	return NewUploadedFile(h, f), nil
}

// Form retrieves value from POST form.
func (request *request) Form(k string) string {
	return request.parent.PostFormValue(k)
}

// Query retrieves value from GET params.
func (request *request) Query(k string) string {
	return request.parent.URL.Query().Get(k)
}

// Param retrieves value from PATH params.
func (request *request) Arg(k string) string {
	if v, ok := request.args[k]; ok {
		return v
	}

	return ""
}

// Input retrieves value with given k name from both Form and Query.
func (request *request) Input(k string) string {
	if v := request.Form(k); v != "" {
		return v
	}

	return request.Query(k)
}

// JSON parse request body as JSON.
func (request *request) JSON(v interface{}) error {
	data, err := ioutil.ReadAll(request.parent.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// IsTLS detects the request is over HTTPS or not.
func (request *request) IsTLS() bool {
	return request.parent.TLS != nil
}

// Header returns original http.Header.
func (request *request) Header() http.Header {
	return request.parent.Header
}

// Method returns request method.
func (request *request) Method() string {
	return request.parent.Method
}

// URI returns RequestURI of incoming requests.
func (request *request) URI() string {
	return request.parent.RequestURI
}

// URL returns URL of incoming requests.
func (request *request) URL() *url.URL {
	return request.parent.URL
}

// Host returns HOST of incoming requests.
func (request *request) Host() string {
	return request.parent.Host
}
