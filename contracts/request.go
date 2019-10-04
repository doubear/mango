package contracts

import (
	"net/http"
	"net/url"
)

// Request represents incoming data from client.
type Request interface {
	Parent() *http.Request
	IP() string
	File(string) (UploadedFile, error)
	Form(string) string
	Query(string) string
	Arg(string) string
	Input(string) string
	JSON(interface{}) error
	IsTLS() bool
	Header() http.Header
	Method() string
	URI() string
	URL() *url.URL
	Host() string
	SetArgs(map[string]string)
	Args() map[string]string
}
