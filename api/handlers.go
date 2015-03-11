package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type jsonMap map[string]interface{}

func registerHandlers() {
	goji.Use(corsMiddleware)

	goji.Get("/", handleIndex)
	goji.Post("/apps/:appSlug/builds", handleBuildCreate)

	goji.Head("/", handleHead)
}

func handleHead(c web.C, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func corsMiddleware(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func sendJson(w http.ResponseWriter, entity interface{}, status int) {
	b, err := json.Marshal(entity)
	if err != nil {
		status = http.StatusInternalServerError
		b = []byte(fmt.Sprintf(`{error:"%s"}`, err.Error()))
	}

	body := string(b[:])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, body)
}

func sendError(w http.ResponseWriter, message string) {
	sendJson(w, jsonMap{"success": false, "error": message}, http.StatusInternalServerError)
}

func handleIndex(c web.C, w http.ResponseWriter, r *http.Request) {
	sendJson(w, jsonMap{"index": true}, http.StatusOK)
}

func handleBuildCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	var app = &App{Slug: c.URLParams["appSlug"]}
	var buildId = time.Now().Format("20060102150405")

	if err := app.Get(); err == nil {
		sendJson(w, jsonMap{"success": true, "build_id": buildId}, http.StatusCreated)
	} else {
		if err.Error() == "The specified key does not exist." {
			sendJson(w, jsonMap{"success": false, "error": "App not found"}, http.StatusNotFound)
		} else {
			sendError(w, err.Error())
		}
	}
}
