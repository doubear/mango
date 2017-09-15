package mango

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/common"
)

func handleResponse(fn common.HandlerFunc) common.MiddleFunc {
	return func(ctx common.Context) {
		code, value := fn(ctx)

		t := reflect.ValueOf(value)

		if code != 0 {
			ctx.Response().SetStatus(code)
		}

		if !t.IsValid() {
			return
		}

		if code == http.StatusPermanentRedirect || code == http.StatusTemporaryRedirect {
			if target, ok := value.(string); ok {
				if target != "" {
					ctx.Response().Redirect(code, target)
				}
			} else {
				panic("calling redirects with an invalid URL target.")
			}

			return
		}

		switch value.(type) {
		case []byte:
			_, err := ctx.Response().Write(value.([]byte))
			if err != nil {
				handleError(ctx, err)
			}
		case string:
			_, err := ctx.Response().WriteString(value.(string))
			if err != nil {
				handleError(ctx, err)
			}
		case *os.File:
			file := value.(*os.File)

			defer file.Close()

			_, err := io.Copy(ctx.Response(), file)
			if err != nil {
				handleError(ctx, err)
			}

			ctx.Response().Header().Set("Content-Disposition", "attachment; filename=\""+file.Name()+"\"")
		case io.Reader:
			_, err := io.Copy(ctx.Response(), value.(io.Reader))
			if err != nil {
				handleError(ctx, err)
			}
		default:
			handleJsonable(ctx, value)
		}

		ctx.Response().Header().Set("Content-Type", http.DetectContentType(ctx.Response().Buffered()))
	}
}

func handleError(ctx common.Context, err error) {
	ctx.Response().SetStatus(http.StatusInternalServerError)
	ctx.Response().Clear()
	logy.W(err.Error())
}

func handleJsonable(ctx common.Context, v interface{}) {
	e := json.NewEncoder(ctx.Response())
	err := e.Encode(v)
	if err != nil {
		handleError(ctx, err)
	}
}
