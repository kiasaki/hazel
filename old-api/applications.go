package main

import (
	"net/http"
)

type ApplicationsEndpoint struct {
	*EndpointBaseWithServer
}

func (e *ApplicationsEndpoint) Index(w http.ResponseWriter, r *http.Request) {
	e.s.tm.RenderPage(w, "applications_index", nil)
}

func (e *ApplicationsEndpoint) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Haii" + GetId(r)))
}
