package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

//header contains JWT encode/decode information.
type header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

func (h *header) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		return ""
	}

	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}
