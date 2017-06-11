package mango

import (
	"encoding/json"
	"io"
	"net/http"
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

		if b, ok := value.([]byte); ok {
			_, err := ctx.W.Write(b)
			if err != nil {
				handleError(ctx, err)
			}

			return
		}

		if b, ok := value.(string); ok {
			_, err := ctx.W.WriteString(b)
			if err != nil {
				handleError(ctx, err)
			}

			return
		}

		if b, ok := value.(io.Reader); ok {
			_, err := io.Copy(ctx.W, b)
			if err != nil {
				handleError(ctx, err)
			}

			return
		}

		handleJsonable(ctx, value)
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
