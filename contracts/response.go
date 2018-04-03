package contracts

import (
	"net/http"
)

//Response 代表即将发送到客户端的响应体
type Response interface {
	Writer() http.ResponseWriter
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
	Redirect(int, string)
	Buffered() []byte
	Send() error
}
