package jwt

import (
	"encoding/base64"
	"strings"
)

func encode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

func decode(s string) []byte {
	if l := len(s) % 4; l > 0 {
		s += strings.Repeat("=", 4-l)
	}

	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return make([]byte, 0)
	}

	return b
}
