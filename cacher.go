package mango

import (
	"time"
)

//Cacher represents Cache specification.
type Cacher interface {
	Open() *Cacher
	Get(id string) (interface{}, error)
	Set(id string, value interface{}, ttl time.Duration) error
	Flush()
	GC()
}
