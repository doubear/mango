package mango

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
)

func handleResponse(fn HandlerFunc) MiddleFunc {
	return func(ctx *Context) {
		code, value := fn(ctx)

		t := reflect.ValueOf(value)

		if code != 0 {
			ctx.W.SetStatus(code)
		}

		if !t.IsValid() {
			return
		}

		if code == http.StatusPermanentRedirect || code == http.StatusTemporaryRedirect {
			if target, ok := value.(string); ok {
				if target != "" {
					ctx.W.Redirect(code, target)
				}
			} else {
				panic("calling redirects with an invalid URL target.")
			}

			return
		}

		switch value.(type) {
		case []byte:
			_, err := ctx.W.Write(value.([]byte))
			if err != nil {
				handleError(ctx, err)
			}
		case string:
			_, err := ctx.W.WriteString(value.(string))
			if err != nil {
				handleError(ctx, err)
			}
		case *os.File:
			file := value.(*os.File)

			defer file.Close()

			_, err := io.Copy(ctx.W, file)
			if err != nil {
				handleError(ctx, err)
			}

			ctx.W.Header().Set("Content-Disposition", "attachment; filename=\""+file.Name()+"\"")
		case io.Reader:
			_, err := io.Copy(ctx.W, value.(io.Reader))
			if err != nil {
				handleError(ctx, err)
			}
		default:
			handleJsonable(ctx, value)
		}

		ctx.W.Header().Set("Content-Type", http.DetectContentType(ctx.W.Buffer()))
	}
}

func handleError(ctx *Context, err error) {
	ctx.W.SetStatus(http.StatusInternalServerError)
	ctx.W.Clear()
	ctx.Logger.Warn(err.Error())
}

func handleJsonable(ctx *Context, v interface{}) {
	e := json.NewEncoder(ctx.W)
	err := e.Encode(v)
	if err != nil {
		handleError(ctx, err)
	}
}
