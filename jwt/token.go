package jwt

import (
	"strings"
	"time"
)

//Token represents JWT token.
type Token struct {
	Header    *header
	Claims    *Claims
	Sign      []byte
	Engine    *Engine
	Audiences []string
}

//Validate checks token validation.
func (t *Token) Validate() bool {

	data := t.Header.String() + "." + t.Claims.String()
	if t.Engine.Algorithm.Verify(t.Engine.Secret, []byte(data), t.Sign) == false {
		panic(ErrTokenInvalid)
	}

	if t.legalAudience() == false {
		panic(ErrTokenInvalid)
	}

	if time.Now().Sub(time.Unix(t.Claims.Exp, 0)) > 0 {
		panic(ErrTokenExpired)
	}

	return true
}

func (t *Token) legalAudience() bool {
	for _, aud := range t.Audiences {
		if strings.Contains(t.Claims.Aud, aud) {
			return true
		}
	}

	return false
}
