package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-mango/mango/contracts"
)

type response struct {
	w      http.ResponseWriter
	io     *bytes.Buffer
	status int
}

//NewResponse create new response instance.
func NewResponse(w http.ResponseWriter) contracts.Response {
	return &response{
		w,
		&bytes.Buffer{},
		http.StatusOK,
	}
}

//Writer returns original http.ResponseWriter
func (r *response) Writer() http.ResponseWriter {
	return r.w
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
	return r.w.Header()
}

//Clear clear buffered data.
func (r *response) Clear() {
	r.io = &bytes.Buffer{}
}

//Size returns total size of response body.
func (r *response) Size() int {
	return r.io.Len()
}

//Status returns status code.
func (r *response) Status() int {
	return r.status
}

//SetStatus reset status of response.
func (r *response) SetStatus(c int) {
	r.status = c
}

//Send sends all buffered data to client.
func (r *response) Send() error {
	r.w.WriteHeader(r.status)
	_, e := r.w.Write(r.io.Bytes())

	if e != nil {
		return e
	}

	r.Clear()

	return nil
}

//SetCookie add a cookie to response header.
func (r *response) SetCookie(c *http.Cookie) {
	http.SetCookie(r.w, c)
}

//DelCookie delete specified cookie.
func (r *response) DelCookie(name string) {
	http.SetCookie(r.w, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})
}

//Redirect recirects to given URL.
func (r *response) Redirect(i int, to string) {
	r.SetStatus(i)
	r.Header().Set("Location", to)
}

//Buffered returns buffered response data.
func (r *response) Buffered() []byte {
	return r.io.Bytes()
}
