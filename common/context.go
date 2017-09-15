package common

//Context represents incoming connection.
type Context interface {
	Request() Request
	Response() Response
	Next()
	Get(string) interface{}
	Set(string, interface{})
	URL(string, map[string]string) string
}
