package contracts

import (
	"net/textproto"
)

//UploadedFile 表示单个文件上传后的具体实例
type UploadedFile interface {
	Header() textproto.MIMEHeader
	Filename() string
	Size() int64
	StoreAs(path string, name string) (string, error)
	Store(path string) (string, error)
}
