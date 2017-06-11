package mango

//Store describes storage interface.
type Store interface {
	Open()
	Close()
}
