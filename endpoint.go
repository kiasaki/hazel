package main

import (
	"net/http"
	"strings"
)

// An Endpoint is destined to handle all requests incomming for a specific
// object, it include all basic CRUD operations
//
// This is a variant on the Ressource type which is more for REST apis this
// Endpoint type is more oriented to web apps
type Endpoint interface {
	Before(http.ResponseWriter, *http.Request)
	After(http.ResponseWriter, *http.Request)
	Prefix() string

	Index(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Store(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Destroy(http.ResponseWriter, *http.Request)
}

type EndpointBase struct {
	prefix string
}

var _ = Endpoint(&EndpointBase{})

func (rs *EndpointBase) Index(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Create(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Store(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Show(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Edit(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Update(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Destroy(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (rs *EndpointBase) Prefix() string {
	return rs.prefix
}

func (rs *EndpointBase) Before(w http.ResponseWriter, r *http.Request) {}
func (rs *EndpointBase) After(w http.ResponseWriter, r *http.Request)  {}

func rHasPrefix(r *http.Request, prefix string) bool {
	return strings.HasPrefix(r.URL.Path, prefix)
}

func rHasSuffix(r *http.Request, suffix string) bool {
	return strings.HasSuffix(r.URL.Path, suffix)
}

// Given a request with a path like
//
// - /users/924924-asd13/edit
// - /users/1
//
// it will extract the middle part based on slashes
//
// Return an empty string if no id can be extracted
func GetId(r *http.Request) string {
	splittedPath := strings.SplitN(r.URL.Path, "/", 4)
	if len(splittedPath) > 2 {
		return splittedPath[2]
	}
	return ""
}

// Calls the right method of the Endpoint based on request method, and
// path format, also calls the Before and After methods
//
// This exists as a separate function and not as ServeHTTP on the Endpoint
// becase it would use EndpointBase implementation of Get and not the child
// one
func NewHandlerForEndpoint(e Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.Before(w, r)
		switch r.Method {
		case "GET":
			if r.URL.Path == e.Prefix() {
				e.Index(w, r)
			} else if r.URL.Path == e.Prefix()+"/create" {
				e.Create(w, r)
			} else if rHasPrefix(r, e.Prefix()) && rHasSuffix(r, "/edit") {
				e.Create(w, r)
			} else {
				e.Show(w, r)
			}
			break
		case "POST":
			e.Store(w, r)
			break
		case "PUT":
		case "PATCH":
			e.Update(w, r)
			break
		case "DELETE":
			e.Destroy(w, r)
			break
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			break
		}
		e.After(w, r)
	})
}

type EndpointBaseWithServer struct {
	*EndpointBase
	s *Server
}

func NewEndpointBaseWithServer(prefix string, s *Server) *EndpointBaseWithServer {
	return &EndpointBaseWithServer{&EndpointBase{prefix}, s}
}
