package contracts

import (
	"net/http"
)

// Response is interface of data will be sent to client.
type Response interface {
	Parent() http.ResponseWriter
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteJSON(interface{}) error
	Header() http.Header
	Clear()
	Size() int
	Status() int
	SetStatus(int)
	SetCookie(*http.Cookie)
	DelCookie(string)
	Redirect(int, string) (int, interface{})
	Buffered() []byte
	Send() error
}
