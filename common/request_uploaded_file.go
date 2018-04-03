package common

import (
	"github.com/go-mango/mango/contracts"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

type uploadedFile struct {
	h *multipart.FileHeader
	f multipart.File
}

//NewUploadedFile 为接收到的文件创建上传文件对象
func NewUploadedFile(h *multipart.FileHeader, f multipart.File) contracts.UploadedFile {
	return &uploadedFile{h, f}
}

//Header 取得原生 textproto.MIMEHeader 实体
func (u *uploadedFile) Header() textproto.MIMEHeader {
	return u.h.Header
}

//Filename 获取上传文件的名称
func (u *uploadedFile) Filename() string {
	return u.h.Filename
}

//Size 获取上传文件的大小
func (u *uploadedFile) Size() int64 {
	return u.h.Size
}

//StoreAs 以指定名称将文件存储到指定路径并返回文件完整路径
func (u *uploadedFile) StoreAs(dst, filename string) (string, error) {
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

//Store 以随即 UUID 为文件名将文件存储到指定路径并返回完整路径
func (u *uploadedFile) Store(dst string) (string, error) {
	return u.StoreAs(dst, uuid.NewV4().String())
}
