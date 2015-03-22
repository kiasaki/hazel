package api

import (
	"github.com/kiasaki/batbelt/rest"
)

type Server struct {
	rest.Server
	Config string
}

func NewServer() Server {
	return Server{rest.NewServer(), ""}
}

func (s *Server) Setup() {
	s.Config = "derp"
}
