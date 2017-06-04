package mango

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Mango struct {
	router   *Router
	middlers []MiddlerFunc
}

func (this *Mango) Use(m MiddlerFunc) {
	this.middlers = append(this.middlers, m)
}

func (this *Mango) Get(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"GET"}, path, fn, middlers)
}

func (this *Mango) Post(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"POST"}, path, fn, middlers)
}

func (this *Mango) Put(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"PUT"}, path, fn, middlers)
}

func (this *Mango) Delete(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"DELETE"}, path, fn, middlers)
}

func (this *Mango) Any(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middlers)
}

func (this *Mango) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params := this.router.Search(r)
	if route == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ms := append(this.middlers, route.MiddlerPool...)

	ctx := newContext(r, w, params, append(ms, MiddleWrapper(route.Handler)))
	ctx.Next()
	ctx.W.flush()
}

//Group create route group with dedicated prefix path.
func (this *Mango) Group(path string, fn GroupFunc, middlers ...MiddlerFunc) {
	path = strings.Trim(path, " /")
	fn(&GroupRouter{
		[]string{path},
		middlers,
		this.router,
	})
}

//Start serve http requests
func (this *Mango) Start(addr string) {
	shouldStop := make(chan os.Signal)
	signal.Notify(shouldStop, os.Interrupt)

	server := &http.Server{
		Addr:    addr,
		Handler: this,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			LogWarn(err.Error())
		}
	}()

	LogInfo("Server is running on " + addr)

	<-shouldStop
	LogWarn("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		LogWarn(err.Error())
	}

	LogInfo("Server stopped gracefully.")
}

//New returns an new Mango instance
func New() *Mango {
	m := &Mango{}

	m.router = &Router{
		make(map[string][]*Route, 0),
		make(map[string][]*Route, 0),
	}

	m.middlers = make([]MiddlerFunc, 0)

	return m
}
