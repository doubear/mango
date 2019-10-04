package contracts

// ThenableFunc uses as a thenable function as a middleware.
type ThenableFunc func(ThenableContext)

// ThenableContext is a context for thenable fucntion.
type ThenableContext interface {
	Context
	Next()
}
