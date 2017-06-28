package jwt

import (
	"encoding/json"
)

//Claims represents JWT claims.
type Claims struct {
	Iss string `json:"iss,omitempty"` //issuer, shows which provider creates the token
	Aud string `json:"aud,omitempty"` //audiences, shows which servers should accept it
	Sub string `json:"sub,omitempty"` //subject, will stores user credential
	Exp int64  `json:"exp,omitempty"` //expire time, references to TTL
	Nbf int64  `json:"nbf,omitempty"`
	Iat int64  `json:"iat,omitempty"` //issued at time
	Jti string `json:"jti,omitempty"`
	Dat int64  `json:"dat,omitempty"` //dead after time, references to RefreshTTL
}

func (c *Claims) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return encode(data)
}
