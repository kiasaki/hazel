package api

import (
	"flag"

	"github.com/kiasaki/batbelt/rest"
)

var appname = "hazel-api"
var version = "0.1.0"

type Server struct {
	rest.Server
	Config Config
}

type Config struct {
	JwtSecret string
}

func NewServer() Server {
	return Server{Server: rest.NewServer(appname, version)}
}

func (s *Server) Setup() {
	// Load config
	s.Config = Config{}
	flag.StringVar(&s.Config.JwtSecret, "jwt-secret", "keyboardcat", "Key to sign JWT tokens")
	flag.Parse()

	// Register services
	loginService := LoginService{s}
	loginService.Register(s.Router)

	// Register endpoints
	s.Register(&ApplicationsEndpoint{s})
}
