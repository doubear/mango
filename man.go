package mango

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

type Mango struct {
	router   *Router
	middlers []MiddlerFunc
}

func (this *Mango) Use(m MiddlerFunc) {
	this.middlers = append(this.middlers, m)
}

//Get register a GET route.
func (this *Mango) Get(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"GET"}, path, fn, middlers)
}

//Post register a POST route.
func (this *Mango) Post(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"POST"}, path, fn, middlers)
}

//Put register a PUT route.
func (this *Mango) Put(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"PUT"}, path, fn, middlers)
}

//Delete register a DELETE route.
func (this *Mango) Delete(path string, fn HandlerFunc, middlers ...MiddlerFunc) {
	this.router.Route([]string{"DELETE"}, path, fn, middlers)
}

//Any register a route without request type limit.
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

func (this *Mango) start(addr string, fn func(*http.Server)) {
	shouldStop := make(chan os.Signal)
	signal.Notify(shouldStop, os.Interrupt)

	server := &http.Server{
		Addr:    addr,
		Handler: this,
	}

	go func() {
		fn(server)
	}()

	defaultLogger.Info("Server is running on " + addr)

	<-shouldStop
	defaultLogger.Warn("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		defaultLogger.Warn(err.Error())
	}

	defaultLogger.Info("Server stopped gracefully.")
}

//Start starts a standard http server.
func (this *Mango) Start(addr string) {
	this.start(addr, func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil {
			defaultLogger.Warn(err.Error())
		}
	})
}

//StartTLS starts a TLS server.
func (this *Mango) StartTLS(addr, certFile, keyFile string) {
	this.start(addr, func(s *http.Server) {
		err := s.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			defaultLogger.Warn(err.Error())
		}
	})
}

//StartAutoTLS starts a TLS server with auto-generated SSL certificate.
//certificates are signed by let's encrypt.
func (this *Mango) StartAutoTLS(addr string, domains ...string) {
	this.start(addr, func(s *http.Server) {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
		}

		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		err := s.ListenAndServeTLS("", "")
		if err != nil {
			defaultLogger.Warn(err.Error())
		}
	})
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
