package mango

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-mango/mango/logger"
)

//Response customized response struct.
type response struct {
	w      http.ResponseWriter
	io     *bytes.Buffer
	status int
	logger *logger.Logger
}

//Write response data to buffer.
func (w *response) Write(b []byte) (int, error) {
	return w.io.Write(b)
}

//WriteString writes string to response body.
func (w *response) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *response) WriteJSON(v interface{}) {
	en := json.NewEncoder(w)
	err := en.Encode(v)

	if err != nil {
		w.logger.Fatal(err.Error())
	}
}

//Header returns http.Header.
func (w *response) Header() http.Header {
	return w.w.Header()
}

//Clear clear buffered data.
func (w *response) Clear() {
	w.io.Reset()
}

//Size returns total size of response body.
func (w *response) Size() int {
	return w.io.Len()
}

//Status returns status code.
func (w *response) Status() int {
	return w.status
}

//SetStatus reset status of response.
func (w *response) SetStatus(c int) {
	w.status = c
}

//Flush all buffered data to client.
func (w *response) flush() {
	w.w.WriteHeader(w.status)
	_, e := w.w.Write(w.io.Bytes())

	if e != nil {
		panic(e.Error())
	}

	w.Clear()
}

//Pusher returns http.Pusher.
func (w *response) Pusher() (http.Pusher, bool) {
	p, ok := w.w.(http.Pusher)

	return p, ok
}

//Push handles http.Pusher within a closure.
func (w *response) Push(fn func(http.Pusher)) bool {
	p, ok := w.Pusher()

	if ok {
		fn(p)
	}

	return ok
}

//Hijacker returns http.Hijacker.
func (w *response) Hijacker() (http.Hijacker, bool) {
	h, ok := w.w.(http.Hijacker)

	return h, ok
}

//Hijack handles http.Hijacker within a closure.
func (w *response) Hijack(fn func(net.Conn, *bufio.ReadWriter)) bool {
	h, ok := w.Hijacker()

	if ok {
		conn, io, err := h.Hijack()

		if err == nil {
			fn(conn, io)
		}
	}

	return ok
}

//SetCookie add a cookie to response header.
func (w *response) SetCookie(c *http.Cookie) {
	http.SetCookie(w.w, c)
}

//DelCookie delete specified cookie.
func (w *response) DelCookie(name string) {
	http.SetCookie(w.w, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})
}

func (w *response) Redirect(i int, to string) {
	w.SetStatus(i)
	w.Header().Set("Location", to)
}
