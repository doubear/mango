package concretes

import (
	"github.com/go-mango/mango/contracts"
)

type session map[string]interface{}

// NewSession create session storage instance.
func NewSession() contracts.Session {
	return session{}
}

func (s session) Set(key string, value interface{}) {
	s[key] = value
}

func (s session) Get(key string) (value interface{}, ok bool) {
	value, ok = s[key]
	return
}
