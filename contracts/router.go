package contracts

// Router interface.
type Router interface {
	Any(string, Callable, ...ThenableFunc)
	Get(string, Callable, ...ThenableFunc)
	Post(string, Callable, ...ThenableFunc)
	Put(string, Callable, ...ThenableFunc)
	Delete(string, Callable, ...ThenableFunc)
	Group(string, func(Router), ...ThenableFunc)
	Use(...ThenableFunc)
	Prefixes() []string
	ThenableStack() []ThenableFunc
	SetThenableStack(...ThenableFunc)
	ToMatch(Request) (Route, map[string]string)
	SetDefaultRoute(Callable)
}
