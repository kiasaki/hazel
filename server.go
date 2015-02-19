package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/kiasaki/batbelt/http/middlewares"
	"github.com/kiasaki/batbelt/http/mm"

	"github.com/kiasaki/mortify/data"
)

type Server struct {
	Cfg Config
	DB  *sql.DB
	mux *http.ServeMux
	tm  *TemplateMap
}

// Creates a new server from config, registers handlers but wait to boot
func NewServer(c Config) *Server {
	s := &Server{Cfg: c}
	s.Setup()
	return s
}

func (s *Server) Setup() {
	s.mux = http.NewServeMux()
	s.DB = data.Connect(s.Cfg.PostgresURL)
	s.registerHandlers()
	s.tm = fillTemplateMap()
}

// Boot the server and starts litening on the configured port
func (s *Server) Run() {
	var stringPort = strconv.Itoa(s.Cfg.Port)
	log.Println("Started listening on port " + stringPort)
	log.Fatal(http.ListenAndServe(":"+stringPort, s.mux))
}

func (s *Server) registerHandlers() {
	mmChain := mm.New(middlewares.LogWithTiming)

	// Add basic auth middleware if configuration specifies user & pass
	u, p := s.Cfg.BasicAuthUser, s.Cfg.BasicAuthPass
	if u != "" && p != "" {
		mmChain = mmChain.Append(middlewares.BasicAuth(u, p))
	}

	s.mux.Handle("/", mmChain.Then(s))
	aEndpoint := ApplicationsEndpoint{NewEndpintBaseWithServer("/apps", s)}
	s.mux.Handle("/apps", mmChain.Then(NewHandlerForEndpoint(&aEndpoint)))
	s.mux.Handle("/apps/", mmChain.Then(NewHandlerForEndpoint(&aEndpoint)))
}

// This handler takes care of 404s and misc paths that don't belong in an
// Endpoint
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/apps", 302)
	} else if r.URL.Path == "/styles.css" {
		RenderStyles(w)
	} else {
		serveNotFound(w, r)
	}
}
