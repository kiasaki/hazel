package main

import (
	"net/http"
)

// A ressource is destined to handle all requests incomming for a specific
// http prefix/ressource
type Ressource interface {
	Before(http.ResponseWriter, *http.Request)
	After(http.ResponseWriter, *http.Request)
	ServeHTTP(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type RessourceBase struct{}

func (rs *RessourceBase) Get(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (rs *RessourceBase) Post(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (rs *RessourceBase) Put(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (rs *RessourceBase) Delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (rs *RessourceBase) Before(w http.ResponseWriter, r *http.Request) {}
func (rs *RessourceBase) After(w http.ResponseWriter, r *http.Request)  {}

// Calls the right method of the ressource based on request method, also calls
// the Before and After methods
func (rs *RessourceBase) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rs.Before(w, r)
	switch r.Method {
	case "GET":
		rs.Get(w, r)
		break
	case "POST":
		rs.Post(w, r)
		break
	case "PUT":
		rs.Put(w, r)
		break
	case "DELETE":
		rs.Delete(w, r)
		break
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		break
	}
	rs.After(w, r)
}

type RessourceBaseWithServer struct {
	RessourceBase
	s *Server
}
