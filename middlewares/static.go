package middlewares

import (
	"net/http"
	"path"
	"strings"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/contracts"

	"os"

	"time"

	"mime"
	"path/filepath"

	"io"
)

//StaticOption configuration of Static middleware.
type StaticOption struct {
	Path string
	Root http.FileSystem
}

//Static serve static assets
func Static(opt StaticOption) contracts.MiddleFunc {

	if opt.Path == "" {
		opt.Path = "/"
		logy.W("StaticOption auto resets Path to /")
	}

	if len(opt.Path) > 1 && opt.Path[0] != '/' {
		opt.Path = "/" + opt.Path
	}

	return func(ctx contracts.Context) {
		fpath := ctx.Request().URL().Path
		if strings.HasPrefix(fpath, opt.Path) {
			fpath = fpath[len(opt.Path):]
			if !strings.HasPrefix(fpath, "/") {
				fpath = "/" + fpath
			}

			fpath = path.Clean(fpath)

			file, err := opt.Root.Open(fpath)
			if err != nil {
				// ctx.Response().SetStatus(resolve(err))
				ctx.Next()
				return
			}

			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				// ctx.Response().SetStatus(resolve(err))
				ctx.Next()
				return
			}

			if stat.IsDir() {
				// ctx.Response().SetStatus(http.StatusForbidden)
				ctx.Next()
				return
			}

			if !stat.ModTime().IsZero() && !stat.ModTime().Equal(time.Unix(0, 0)) {
				ctx.Response().Header().Add("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))
			}

			if _, ok := ctx.Response().Header()["Content-Type"]; !ok {
				m := mime.TypeByExtension(filepath.Ext(stat.Name()))
				if m == "" {
					m = "application/octet-stream"
				}

				ctx.Response().Header().Add("Content-Type", m)
			}

			_, err = io.Copy(ctx.Response(), file)
			if err != nil {
				ctx.Response().Clear()
				ctx.Response().SetStatus(http.StatusInternalServerError)
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
