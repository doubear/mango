package mango

import (
	pcontext "context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-mango/logy"

	"golang.org/x/crypto/acme/autocert"
)

//Plugin an plugin for mango.
type Plugin func(*Mango)

//Mango main struct.
type Mango struct {
	router   *router
	middles  []MiddleFunc
	notFound HandlerFunc
	cacher   Cacher
}

//NewContext create new Context instance
func (m *Mango) newContext(r *http.Request, w http.ResponseWriter, ps map[string]string, ms []MiddleFunc) *context {
	return &context{
		newRequest(r),
		newResponse(w),
		m.cacher,
		ps,
		ms,
		map[string]interface{}{},
	}
}

//SetCacher sets cache provider.
func (m *Mango) SetCacher(c Cacher) {
	m.cacher = c
}

//Use appends middleware function to built-in stack,
//or load plugin to framework.
func (m *Mango) Use(fn interface{}) {
	switch fn.(type) {
	case MiddleFunc:
		m.middles = append(m.middles, fn.(MiddleFunc))
	case Plugin:
		fn.(Plugin)(m)
	default:
		logy.E("use an invalid value")
	}
}

//NotFound set customized not found error handler.
func (m *Mango) NotFound(fn HandlerFunc) {
	m.notFound = fn
}

//Get register a GET route.
func (m *Mango) Get(path string, fn HandlerFunc, middles ...MiddleFunc) {
	m.router.route([]string{"GET"}, path, fn, middles)
}

//Post register a POST route.
func (m *Mango) Post(path string, fn HandlerFunc, middles ...MiddleFunc) {
	m.router.route([]string{"POST"}, path, fn, middles)
}

//Put register a PUT route.
func (m *Mango) Put(path string, fn HandlerFunc, middles ...MiddleFunc) {
	m.router.route([]string{"PUT"}, path, fn, middles)
}

//Delete register a DELETE route.
func (m *Mango) Delete(path string, fn HandlerFunc, middles ...MiddleFunc) {
	m.router.route([]string{"DELETE"}, path, fn, middles)
}

//Any register a route without request type limit.
func (m *Mango) Any(path string, fn HandlerFunc, middles ...MiddleFunc) {
	m.router.route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middles)
}

func (m *Mango) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rt *route
	found, params := m.router.search(r)
	if found == nil {
		rt = &route{
			method:     "*",
			path:       "/",
			handler:    m.notFound,
			middlePool: make([]MiddleFunc, 0),
		}
	} else {
		rt = found
	}

	ms := append(m.middles, rt.middlePool...)

	ctx := m.newContext(r, w, params, append(ms, handleResponse(rt.handler)))
	ctx.Next()
	ctx.Response().Send()
}

//Group create route group with dedicated prefix path.
func (m *Mango) Group(path string, fn GroupFunc, middles ...MiddleFunc) {
	path = strings.Trim(path, " /")
	fn(&GroupRouter{
		[]string{path},
		middles,
		m.router,
	})
}

func (m *Mango) start(addr string, fn func(*http.Server)) {
	shouldStop := make(chan os.Signal)
	signal.Notify(shouldStop, os.Interrupt, os.Kill)

	server := &http.Server{
		Addr:    addr,
		Handler: m,
	}

	go func() {
		fn(server)
	}()

	logy.I("Server is running on " + addr)

	<-shouldStop
	logy.W("Server is shutting down...")

	ctx, cancel := pcontext.WithTimeout(pcontext.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logy.W(err.Error())
	}

	logy.I("Server stopped gracefully.")
}

//Start starts a standard http server.
func (m *Mango) Start(addr string) {
	m.start(addr, func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil {
			logy.W(err.Error())
		}
	})
}

//StartTLS starts a TLS server.
func (m *Mango) StartTLS(addr, certFile, keyFile string) {
	m.start(addr, func(s *http.Server) {
		err := s.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			logy.W(err.Error())
		}
	})
}

//StartAutoTLS starts a TLS server with auto-generated SSL certificate.
//certificates are signed by let's encrypt.
func (m *Mango) StartAutoTLS(addr string, caStore autocert.Cache, domains ...string) {
	m.start(addr, func(s *http.Server) {
		c := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      caStore,
		}

		s.TLSConfig = &tls.Config{GetCertificate: c.GetCertificate}

		err := s.ListenAndServeTLS("", "")
		if err != nil {
			logy.W(err.Error())
		}
	})
}

//New returns an new Mango instance
func New() *Mango {
	m := &Mango{}

	m.router = &router{
		make(map[string][]*route, 0),
		make(map[string][]*route, 0),
	}

	m.notFound = func(ctx Context) (int, interface{}) {
		ctx.Response().SetStatus(http.StatusNotFound)
		return 0, nil
	}

	m.middles = make([]MiddleFunc, 0)

	return m
}

//Default returns an Mango instance that uses few middlewares.
func Default() *Mango {
	m := New()
	m.Use(Record())
	m.Use(Recovery())
	return m
}
