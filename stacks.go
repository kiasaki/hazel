package main

import (
	"net/http"
)

type StacksEndpoint struct {
	*EndpointBaseWithServer
}

func (e *StacksEndpoint) Index(w http.ResponseWriter, r *http.Request) {
	e.s.tm.RenderPage(w, "stacks_index", nil)
}

func (e *StacksEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	e.s.tm.RenderPage(w, "stacks_create", nil)
}

func (e *StacksEndpoint) Store(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
}
