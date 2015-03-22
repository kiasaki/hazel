package api

import (
	"github.com/kiasaki/batbelt/rest"
)

var appname = "hazel-api"
var version = "0.1.0"

type Server struct {
	rest.Server
	Config string
}

func NewServer() Server {
	return Server{Server: rest.NewServer(appname, version)}
}

func (s *Server) Setup() {
	// Load config
	s.Config = "derp"

	// Register endpoints
	s.Register(&ApplicationsEndpoint{s})
}
