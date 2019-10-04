package mango

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/concretes"
	"github.com/go-mango/mango/contracts"
	"github.com/go-mango/mango/middlewares"

	"golang.org/x/crypto/acme/autocert"
)

type mango struct {
	router    contracts.Router
	thenStack []contracts.ThenableFunc
	cache     contracts.Cachable
	events    map[string][]func()
}

func (m *mango) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := concretes.NewRequest(r)
	response := concretes.NewResponse(w)

	route, params := m.router.ToMatch(request)
	request.SetArgs(params)
	thenStack := append(m.thenStack, route.ThenStack()...)

	ctx := concretes.NewContext(
		request,
		response,
		m.cache,
		thenStack,
		route,
	)

	ctx.Next()
	ctx.Response().Send()
}

//SetCachable sets cache provider.
func (m *mango) SetCachable(cache contracts.Cachable) {
	m.cache = cache
}

//Use appends contracts.ThenableFunc to built-in stack.
func (m *mango) Use(next contracts.ThenableFunc) {
	m.thenStack = append(m.thenStack, next)
}

//SetDefaultRoute set customized not found error handler.
func (m *mango) SetDefaultRoute(fn contracts.Callable) {
	m.router.SetDefaultRoute(fn)
}

//Get register a GET route.
func (m *mango) Get(path string, fn contracts.Callable, thenStack ...contracts.ThenableFunc) {
	m.router.Get(path, fn, thenStack...)
}

//Post register a POST route.
func (m *mango) Post(path string, fn contracts.Callable, thenStack ...contracts.ThenableFunc) {
	m.router.Post(path, fn, thenStack...)
}

//Put register a PUT route.
func (m *mango) Put(path string, fn contracts.Callable, thenStack ...contracts.ThenableFunc) {
	m.router.Put(path, fn, thenStack...)
}

//Delete register a DELETE route.
func (m *mango) Delete(path string, fn contracts.Callable, thenStack ...contracts.ThenableFunc) {
	m.router.Delete(path, fn, thenStack...)
}

//Any register a route without request type limit.
func (m *mango) Any(path string, fn contracts.Callable, thenStack ...contracts.ThenableFunc) {
	m.router.Any(path, fn, thenStack...)
}

//Group create route group with dedicated prefix path.
func (m *mango) Group(path string, fn func(contracts.Router), thenStack ...contracts.ThenableFunc) {
	m.router.Group(path, fn, thenStack...)
}

func (m *mango) start(addr string, fn func(*http.Server)) {
	shouldStop := make(chan os.Signal)
	signal.Notify(shouldStop, os.Interrupt, os.Kill)

	server := &http.Server{
		Addr:    addr,
		Handler: m,
	}

	m.emit("starting")
	go func() {
		fn(server)
	}()

	m.emit("started")

	logy.Std().Info("Server is running on", addr)

	<-shouldStop
	logy.Std().Warn("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logy.Std().Warn(err.Error())
	}

	logy.Std().Info("Server is stopping daemon tasks...")
	m.emit("shutdown")

	logy.Std().Info("Server stopped gracefully.")
}

//Start starts a standard http server.
func (m *mango) Start(addr string) {
	m.start(addr, func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil {
			logy.Std().Error(err.Error())
		}
	})
}

//StartTLS starts a TLS server.
func (m *mango) StartTLS(addr, certFile, keyFile string) {
	m.start(addr, func(s *http.Server) {
		err := s.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			logy.Std().Error(err.Error())
		}
	})
}

// StartAutoTLS starts a TLS server with auto-generated SSL certificate.
// certificates are signed by let's encrypt.
func (m *mango) StartAutoTLS(addr string, caStore autocert.Cache, domains ...string) {
	m.start(addr, func(s *http.Server) {
		c := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      caStore,
		}

		s.TLSConfig = &tls.Config{GetCertificate: c.GetCertificate}

		err := s.ListenAndServeTLS("", "")
		if err != nil {
			logy.Std().Error(err.Error())
		}
	})
}

// On regists callback function while specified event is emitted.
func (m *mango) On(event string, fn func()) {
	if _, ok := m.events[event]; ok == false {
		m.events[event] = make([]func(), 0)
	}

	m.events[event] = append(m.events[event], fn)
}

func (m *mango) emit(event string) {
	if fns, ok := m.events[event]; ok {
		for _, fn := range fns {
			fn()
		}
	}
}

// New returns an new Mango instance
func New() contracts.Mango {
	m := &mango{
		concretes.NewRouter(),
		[]contracts.ThenableFunc{},
		concretes.NewMemoryCache(15 * time.Minute),
		map[string][]func(){},
	}

	m.SetDefaultRoute(func(ctx contracts.Context) (int, interface{}) {
		return 404, "page not found"
	})

	m.thenStack = []contracts.ThenableFunc{}

	return m
}

// Default returns an Mango instance that uses few middlewares.
func Default() contracts.Mango {
	m := New()
	m.Use(middlewares.Record())
	m.Use(middlewares.Recovery())
	return m
}
