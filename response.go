package mango

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
)

//Response customized response struct.
type response struct {
	w      http.ResponseWriter
	io     *bytes.Buffer
	status int
}

//Write response data to buffer.
func (this *response) Write(b []byte) (int, error) {
	return this.io.Write(b)
}

//WriteString writes string to response body.
func (this *response) WriteString(s string) (int, error) {
	return this.Write([]byte(s))
}

func (this *response) WriteJSON(v interface{}) {
	en := json.NewEncoder(this)
	err := en.Encode(v)

	if err != nil {
		LogFatal(err.Error())
	}
}

//Header returns http.Header.
func (this *response) Header() http.Header {
	return this.w.Header()
}

//Clear clear buffered data.
func (this *response) Clear() {
	this.io.Reset()
}

//Size returns total size of response body.
func (this *response) Size() int {
	return this.io.Len()
}

//Status returns status code.
func (this *response) Status() int {
	return this.status
}

//SetStatus reset status of response.
func (this *response) SetStatus(c int) {
	this.status = c
}

//Flush all buffered data to client.
func (this *response) flush() {
	this.w.WriteHeader(this.status)
	_, e := this.w.Write(this.io.Bytes())

	if e != nil {
		panic(e.Error())
	}

	this.Clear()
}

//Pusher returns http.Pusher.
func (this *response) Pusher() (http.Pusher, bool) {
	p, ok := this.w.(http.Pusher)

	return p, ok
}

//Push handles http.Pusher within a closure.
func (this *response) Push(fn func(http.Pusher)) bool {
	p, ok := this.Pusher()

	if ok {
		fn(p)
	}

	return ok
}

//Hijacker returns http.Hijacker.
func (this *response) Hijacker() (http.Hijacker, bool) {
	h, ok := this.w.(http.Hijacker)

	return h, ok
}

//Hijack handles http.Hijacker within a closure.
func (this *response) Hijack(fn func(net.Conn, *bufio.ReadWriter)) bool {
	h, ok := this.Hijacker()

	if ok {
		conn, io, err := h.Hijack()

		if err == nil {
			fn(conn, io)
		}
	}

	return ok
}

func newResponse(w http.ResponseWriter) *response {
	return &response{
		w,
		&bytes.Buffer{},
		http.StatusOK,
	}
}
