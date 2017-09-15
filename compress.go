package mango

import (
	"compress/gzip"
	"strings"
)

//Compress compress response data.
func Compress() MiddleFunc {
	return func(ctx Context) {
		ctx.Next() //continues to execute middlewares.

		if ctx.Response().Size() == 0 {
			return
		}

		accept := ctx.Request().Header().Get("Accept-Encoding")
		if strings.Contains(accept, "gzip") {
			data := ctx.Response().Buffered()
			ctx.Response().Clear()

			w, err := gzip.NewWriterLevel(ctx.Response(), 6)
			if err != nil {
				ctx.Response().Clear()
				ctx.Response().Write(data)
				return
			}

			_, err = w.Write(data)
			w.Close()
			if err != nil {
				ctx.Response().Clear()
				ctx.Response().Write(data)
				return
			}

			ctx.Response().Header().Set("Content-Encoding", "gzip")
			ctx.Response().Header().Set("Vary", "Accept-Encoding")
		}
	}
}
