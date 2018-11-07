package contracts

import (
	"time"
)

//Cacher represents Cache specification.
//
/*
	m.Use(cache.Memory())

	ctx.C.Get("xxx")
*/
type Cacher interface {
	Get(id string) interface{}
	Set(id string, value interface{}, ttl time.Duration)
	Del(id string)
	Push(id string, value interface{})
	Pop(id string) interface{}
	Flush()
	GC()
}
