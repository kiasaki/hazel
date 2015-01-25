package main

import (
	"log"
	"net/http"

	"github.com/kiasaki/hazel/api"
	"github.com/kiasaki/hazel/data"
)

func main() {
	log.Println("Starting Server")

	http.Handle("/api/", api.Mux())
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui/"))))
	http.HandleFunc("/", handleRoot)

	log.Println("Ensuring indexes")
	data.Index(data.Database())

	log.Println("Listening on 4411")
	http.ListenAndServe(":4411", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/", 301)
}
