package contracts

//Context represents incoming connection.
type Context interface {
	Request() Request
	Response() Response
	Auth() Authenable
	URL(string, map[string]string) string
	Cache() Cachable
	Session() Session
}
