package jwt

import (
	"crypto"
	"crypto/hmac"
)

//Crypto represents crypto methods.
type Crypto struct {
	Name string
	Hash crypto.Hash
}

var (
	//HS256 references to sha256.
	HS256 = Crypto{"HS256", crypto.SHA256}
)

//Sign makes hash value of given data.
func (c *Crypto) Sign(k, s string) string {
	h := hmac.New(c.Hash.New, []byte(k))
	h.Write([]byte(s))
	return encode(h.Sum(nil))
}

//Verify creates hash and validate it.
func (c *Crypto) Verify(k string, data, sign []byte) bool {
	h := hmac.New(c.Hash.New, []byte(k))
	h.Write(data)

	return hmac.Equal(sign, h.Sum(nil))
}
