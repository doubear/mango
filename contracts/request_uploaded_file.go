package contracts

import (
	"net/textproto"
)

// UploadedFile is interface of uploaded files handler.
type UploadedFile interface {
	Header() textproto.MIMEHeader
	Filename() string
	Size() int64
	StoreAs(path string, name string) (string, error)
	Store(path string) (string, error)
}
