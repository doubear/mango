package common

import (
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

//UploadedFile represent uploaded file.
type UploadedFile struct {
	h *multipart.FileHeader
	f multipart.File
}

//Header returns original textproto.MIMEHeader.
func (u *UploadedFile) Header() textproto.MIMEHeader {
	return u.h.Header
}

//Filename returns the name of uploaded file.
func (u *UploadedFile) Filename() string {
	return u.h.Filename
}

//Size returns file size of uploaded file.
func (u *UploadedFile) Size() int64 {
	return u.h.Size
}

//StoreAs stores uploaded file to dst with filename,
//then returns it file path.
func (u *UploadedFile) StoreAs(dst, filename string) (string, error) {
	storedAt := filepath.Join(dst, filename)
	f, err := os.Create(storedAt)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, u.f)
	if err != nil {
		return "", err
	}

	return storedAt, nil
}

//Store stores uploaded file to dst with randomized filename,
//then returns it file path.
func (u *UploadedFile) Store(dst string) (string, error) {
	return u.StoreAs(dst, uuid.NewV4().String())
}
