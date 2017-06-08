package mango

import (
	"net/http"
	"path"
	"strings"

	"os"

	"time"

	"mime"
	"path/filepath"

	"io"

	"github.com/go-mango/mango/logger"
)

//StaticOption configuration of Static middleware.
type StaticOption struct {
	Path string
	Root http.FileSystem
}

//Static serve static assets
func Static(opt StaticOption) MiddleFunc {

	if opt.Path == "" {
		opt.Path = "/"
		logger.NewLogger().Warn("StaticOption auto resets Path to /")
	}

	if len(opt.Path) > 1 && opt.Path[0] != '/' {
		opt.Path = "/" + opt.Path
	}

	return func(ctx *Context) {
		fpath := ctx.R.URL.Path
		if strings.HasPrefix(fpath, opt.Path) {
			fpath = fpath[len(opt.Path):]
			if !strings.HasPrefix(fpath, "/") {
				fpath = "/" + fpath
			}

			fpath = path.Clean(fpath)

			file, err := opt.Root.Open(fpath)
			if err != nil {
				ctx.W.SetStatus(resolve(err))
				return
			}

			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				ctx.W.SetStatus(resolve(err))
				return
			}

			if stat.IsDir() {
				ctx.W.SetStatus(http.StatusForbidden)
				return
			}

			if !stat.ModTime().IsZero() && !stat.ModTime().Equal(time.Unix(0, 0)) {
				ctx.W.Header().Add("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))
			}

			if _, ok := ctx.W.Header()["Content-Type"]; !ok {
				m := mime.TypeByExtension(filepath.Ext(stat.Name()))
				if m == "" {
					m = "application/octet-stream"
				}

				ctx.W.Header().Add("Content-Type", m)
			}

			_, err = io.Copy(ctx.W, file)
			if err != nil {
				ctx.W.Clear()
				ctx.W.SetStatus(http.StatusInternalServerError)
			}

			return
		}

		ctx.Next()
	}
}

func resolve(e error) int {
	if os.IsNotExist(e) {
		return http.StatusNotFound
	}

	if os.IsPermission(e) {
		return http.StatusForbidden
	}

	return http.StatusInternalServerError
}
