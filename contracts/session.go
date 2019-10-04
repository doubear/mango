package contracts

// Session interface.
type Session interface {
	Set(string, interface{})
	Get(string) (interface{}, bool)
}
