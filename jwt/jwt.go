package jwt

import (
	"encoding/json"
	"net/http"
	"strings"

	"time"

	"github.com/go-mango/mango"
)

var (
	defaultErrorHandler = func(e error, ctx *mango.Context) {
		ctx.W.Clear()
		ctx.W.SetStatus(http.StatusUnauthorized)
	}

	defaultTTL        = 2 * time.Hour
	defaultRefreshTTL = 7 * 24 * time.Hour
)

//Author returns subject of jwt claims.
type Author interface {
	GetSubject() string
}

//Engine represents json-web-token.
type Engine struct {
	Secret     string
	TTL        time.Duration
	RefreshTTL time.Duration
	Algorithm  Crypto
	header     *header
	onError    func(error, *mango.Context)
}

//Error set unauthorized handler.
func (e *Engine) Error(fn func(error, *mango.Context)) {
	e.onError = fn
}

//Sign creates hash value of given claims.
func (e *Engine) Sign(c Claims) string {
	c.Iat = time.Now().Unix()
	c.Exp = time.Now().Add(e.TTL).Unix()
	c.Dat = time.Now().Add(e.RefreshTTL).Unix()

	data := e.header.String() + "." + c.String()
	data = data + "." + e.Algorithm.Sign(e.Secret, data)

	return data
}

//ParseToken parse jwt token.
func (e *Engine) ParseToken(t string) *Token {
	parts := strings.SplitN(t, ".", 3)
	if len(parts) != 3 {
		return nil
	}

	h := &header{}
	err := json.Unmarshal(decode(parts[0]), h)
	if err != nil {
		return nil
	}

	c := &Claims{}
	err = json.Unmarshal(decode(parts[1]), c)
	if err != nil {
		return nil
	}

	return &Token{h, c, decode(parts[2]), e, []string{}}
}

//Auth returns jwt authorization middleware of mango.
func (e *Engine) Auth(audiences ...string) mango.MiddleFunc {
	return func(ctx *mango.Context) {
		defer func() {
			if rev := recover(); rev != nil {
				e.onError(rev.(error), ctx)
			}
		}()

		if token := ctx.R.Header.Get("Authorization"); token != "" {
			if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
				t := e.ParseToken(token)
				t.Audiences = audiences
				if t != nil && t.Validate() {
					ctx.Set("jwt.sub", t.Claims.Sub)
					ctx.Set("jwt.aud", t.Claims.Aud)
					ctx.Set("jwt.iss", t.Claims.Iss)
					ctx.Next()
					return
				}
			}
		}

		panic(ErrTokenLost)
	}
}

//New creates new jwt instance.
//
/*
	e := jwt.New(jwt.HS256, "secret key")
	e.Sign(jwt.Claims{})
*/
func New(alg Crypto, secret string) *Engine {
	return &Engine{
		secret,
		defaultTTL,
		defaultRefreshTTL,
		alg,
		&header{"JWT", alg.Name},
		defaultErrorHandler,
	}
}
