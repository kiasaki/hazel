package main

import (
	"net/http"
)

func serveNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404 - Page not found"))
}
