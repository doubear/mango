package common

//HandlerFunc use to handle incoming requests.
type HandlerFunc func(Context) (int, interface{})
