package api

import (
	"net/http"

	"github.com/kiasaki/batbelt/rest"
)

type ApplicationsEndpoint struct {
	s *Server
}

func (e *ApplicationsEndpoint) Path() string {
	return "/applications"
}

func (e *ApplicationsEndpoint) GET(w http.ResponseWriter, r *http.Request) {
	rest.SetOKResponse(w, rest.J{"applications": []string{}})
}
