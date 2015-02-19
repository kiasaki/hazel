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
}

// Router the request to the right ressource based on the url path prefix
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Write([]byte("Hai"))
	} else {
		serveNotFound(w, r)
	}
}
