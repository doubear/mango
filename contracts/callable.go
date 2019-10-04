package contracts

// Callable use to handle incoming requests.
type Callable func(Context) (int, interface{})
