package contracts

import (
	"regexp"
)

// Route is part of Router.
type Route interface {
	Method() string
	Path() string
	SetPath(string)
	Pathable() *regexp.Regexp
	Callable() Callable
	ThenStack() []ThenableFunc
	IsStatic() bool
}
