package storage

import (
	"time"
)

//Adapter storage driver adapter.
type Adapter interface {
	Put(i string, v interface{}, t time.Duration) error
	Get(i string) (interface{}, error)
	Del(i string)
	Has(i string) bool
	Clear()
	GC()
	Open()
	Close()
}
