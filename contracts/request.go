package contracts

import (
	"net/url"
	"net/http"
)

//Request 代表单个入站请求
type Request interface {
	Parent() *http.Request
	IP() string
	File(string) (UploadedFile, error)
	Form(string) string
	Query(string) string
	Param(string) string
	Input(string) string
	JSON(interface{}) error
	IsTLS() bool
	Header() http.Header
	Method() string
	URI() string
	URL() *url.URL
	Host() string
}
