package mango

import (
	pcontext "context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/crypto/acme"

	"github.com/go-mango/logy"
	"github.com/go-mango/mango/contracts"
	mhttp "github.com/go-mango/mango/http"
	"github.com/go-mango/mango/middlewares"

	"golang.org/x/crypto/acme/autocert"
)

const (
	leStagingDirectory = "https://acme-staging.api.letsencrypt.org/directory"
	stagingMode        = iota
	prodMode
)

//Plugin an plugin for mango.
type Plugin func(*Mango)

//Mango main struct.
type Mango struct {
	router   *router
	middles  []contracts.MiddleFunc
	notFound contracts.HandlerFunc
	cacher   contracts.Cacher
	mode     int
}

//NewContext create new Context instance
func (m *Mango) newContext(r *http.Request, w http.ResponseWriter, ps map[string]string, ms []contracts.MiddleFunc) *context {
	return &context{
		mhttp.NewRequest(r, ps),
		mhttp.NewResponse(w),
		m.cacher,
		ms,
		map[string]interface{}{},
	}
}

//SetCacher sets cache provider.
func (m *Mango) SetCacher(c contracts.Cacher) {
	m.cacher = c
}

//Use appends middleware function to built-in stack,
//or load plugin to framework.
func (m *Mango) Use(fn interface{}) {
	switch fn.(type) {
	case contracts.MiddleFunc:
		m.middles = append(m.middles, fn.(contracts.MiddleFunc))
	case Plugin:
		fn.(Plugin)(m)
	default:
		logy.Std().Error("use an invalid value")
	}
}

//NotFound set customized not found error handler.
func (m *Mango) NotFound(fn contracts.HandlerFunc) {
	m.notFound = fn
}

//Get register a GET route.
func (m *Mango) Get(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	m.router.route([]string{"GET"}, path, fn, middles)
}

//Post register a POST route.
func (m *Mango) Post(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	m.router.route([]string{"POST"}, path, fn, middles)
}

//Put register a PUT route.
func (m *Mango) Put(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	m.router.route([]string{"PUT"}, path, fn, middles)
}

//Delete register a DELETE route.
func (m *Mango) Delete(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	m.router.route([]string{"DELETE"}, path, fn, middles)
}

//Any register a route without request type limit.
func (m *Mango) Any(path string, fn contracts.HandlerFunc, middles ...contracts.MiddleFunc) {
	m.router.route([]string{"GET", "POST", "PUT", "DELETE"}, path, fn, middles)
}

func (m *Mango) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	found, params := m.router.search(r)
	if found == nil {
		found = &route{
			method:     "*",
			path:       "/",
			handler:    m.notFound,
			middlePool: make([]contracts.MiddleFunc, 0),
		}
	}

	ms := append(m.middles, found.middlePool...)

	ctx := m.newContext(r, w, params, append(ms, handleResponse(found.handler)))
	ctx.Next()
	ctx.Response().Send()
}

//Group create route group with dedicated prefix path.
func (m *Mango) Group(path string, fn GroupFunc, middles ...contracts.MiddleFunc) {
	path = strings.Trim(path, " /")
	fn(&GroupRouter{
		[]string{path},
		middles,
		m.router,
	})
}

//RunInStagingMode sets staging mode in runtime.
func (m *Mango) RunInStagingMode() {
	m.mode = stagingMode
}

//RunInProdMode sets production mode in runtime.
func (m *Mango) RunInProdMode() {
	m.mode = prodMode
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

	logy.Std().Info("Server is running on " + addr)

	<-shouldStop
	logy.Std().Warn("Server is shutting down...")

	ctx, cancel := pcontext.WithTimeout(pcontext.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logy.Std().Warn(err.Error())
	}

	logy.Std().Info("Server stopped gracefully.")
}

//Start starts a standard http server.
func (m *Mango) Start(addr string) {
	m.start(addr, func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil {
			logy.Std().Error(err.Error())
		}
	})
}

//StartTLS starts a TLS server.
func (m *Mango) StartTLS(addr, certFile, keyFile string) {
	m.start(addr, func(s *http.Server) {
		err := s.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			logy.Std().Error(err.Error())
		}
	})
}

//StartAutoTLS starts a TLS server with auto-generated SSL certificate.
//certificates are signed by let's encrypt.
func (m *Mango) StartAutoTLS(addr string, caStore autocert.Cache, domains ...string) {
	var caClient *acme.Client
	if m.mode == stagingMode {
		caClient = &acme.Client{DirectoryURL: leStagingDirectory}
	}

	m.start(addr, func(s *http.Server) {
		c := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      caStore,
			Client:     caClient,
		}

		s.TLSConfig = &tls.Config{GetCertificate: c.GetCertificate}

		err := s.ListenAndServeTLS("", "")
		if err != nil {
			logy.Std().Error(err.Error())
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

	m.notFound = func(ctx contracts.Context) (int, interface{}) {
		ctx.Response().SetStatus(http.StatusNotFound)
		return 0, nil
	}

	m.middles = make([]contracts.MiddleFunc, 0)

	return m
}

//Default returns an Mango instance that uses few middlewares.
func Default() *Mango {
	m := New()
	m.Use(middlewares.Record())
	m.Use(Recovery())
	return m
}
