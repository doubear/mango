package mango

import (
	"compress/gzip"
	"strings"
)

//Compress compress response data.
func Compress() MiddleFunc {
	return func(ctx *Context) {
		ctx.Next() //continues to execute middlewares.

		if ctx.W.Size() == 0 {
			return
		}

		accept := ctx.R.Header.Get("Accept-Encoding")
		if strings.Contains(accept, "gzip") {
			data := ctx.W.Buffer()
			ctx.W.Clear()

			w, err := gzip.NewWriterLevel(ctx.W, 6)
			if err != nil {
				ctx.W.Clear()
				ctx.W.Write(data)
				return
			}

			_, err = w.Write(data)
			w.Close()
			if err != nil {
				ctx.W.Clear()
				ctx.W.Write(data)
				return
			}

			ctx.W.Header().Set("Content-Encoding", "gzip")
			ctx.W.Header().Set("Vary", "Accept-Encoding")
		}
	}
}
