package api

import (
	"flag"

	"github.com/kiasaki/batbelt/rest"
	"github.com/kiasaki/hazel/api/data"
)

var appname = "hazel-api"
var version = "0.1.0"

type Server struct {
	rest.Server
	Config Config
	DB     *data.Database
}

type Config struct {
	JwtSecret string
	DBFile    string
}

func NewServer() Server {
	return Server{Server: rest.NewServer(appname, version)}
}

func (s *Server) Setup() {
	// Load config
	s.Config = Config{}
	flag.StringVar(&s.Config.JwtSecret, "jwt-secret", "keyboardcat", "Key to sign JWT tokens")
	flag.StringVar(&s.Config.DBFile, "db-file", "hazel.db", "BoldDB file location")
	flag.Parse()

	// Setup database
	s.DB = data.NewDatabase(s.Config.DBFile)

	// Register services
	loginService := LoginService{s}
	loginService.Register(s.Router)

	// Register endpoints
	s.Register(&ApplicationsEndpoint{s})
}

func (s *Server) Run() {
	err := s.DB.Open()
	if err != nil {
		s.Logger.Fatal(err)
	}
	defer s.DB.Close()

	s.Server.Run()
}
