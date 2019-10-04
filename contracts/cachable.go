package contracts

import (
	"time"
)

//Cachable represents Cachable specification.
type Cachable interface {
	Get(id string) interface{}
	Set(id string, value interface{}, ttl time.Duration)
	Del(id string)
	Push(id string, value interface{})
	Pop(id string) interface{}
	Flush()
	GC()
}
