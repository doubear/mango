package concretes

import (
	"github.com/go-mango/mango/contracts"
)

type auth struct {
	userID string
}

// newAuth creates authenable instance.
func newAuth() contracts.Authenable {
	return new(auth)
}

func (auth *auth) SetUserID(id string) {
	auth.userID = id
}

func (auth *auth) UserID() string {
	return auth.userID
}
