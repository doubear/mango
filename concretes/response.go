package concretes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-mango/mango/contracts"
)

type response struct {
	parent http.ResponseWriter
	io     *bytes.Buffer
	status int
}

// NewResponse create new response instance.
func NewResponse(parent http.ResponseWriter) contracts.Response {
	return &response{
		parent,
		&bytes.Buffer{},
		http.StatusOK,
	}
}

//Parent returns original http.ResponseWriter
func (r *response) Parent() http.ResponseWriter {
	return r.parent
}

//Write response data to buffer.
func (r *response) Write(b []byte) (int, error) {
	return r.io.Write(b)
}

//WriteString writes string to response body.
func (r *response) WriteString(s string) (int, error) {
	return r.Write([]byte(s))
}

func (r *response) WriteJSON(v interface{}) error {
	return json.NewEncoder(r).Encode(v)
}

//Header returns http.Header.
func (r *response) Header() http.Header {
	return r.parent.Header()
}

//Clear clear buffered data.
func (r *response) Clear() {
	r.io = &bytes.Buffer{}
}

//Size returns total size of response body.
func (r *response) Size() int {
	return r.io.Len()
}

// Status returns status code.
func (r *response) Status() int {
	return r.status
}

// SetStatus reset status of response.
func (r *response) SetStatus(c int) {
	r.status = c
}

// Send sends all buffered data to client.
func (r *response) Send() error {
	r.parent.WriteHeader(r.status)
	_, e := io.Copy(r.parent, r.io)

	if e != nil {
		return e
	}

	r.Clear()

	return nil
}

// SetCookie add a cookie to response header.
func (r *response) SetCookie(c *http.Cookie) {
	http.SetCookie(r.parent, c)
}

// DelCookie delete specified cookie.
func (r *response) DelCookie(name string) {
	http.SetCookie(r.parent, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})
}

// Redirect recirects to given URL.
func (r *response) Redirect(status int, to string) (int, interface{}) {
	r.SetStatus(status)
	r.Header().Set("Location", to)
	return 0, nil
}

//Buffered returns buffered response data.
func (r *response) Buffered() []byte {
	return r.io.Bytes()
}
