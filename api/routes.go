package api

import (
	"net/http"

	"github.com/kiasaki/hazel/helper"
)

var APIIntance *helper.API

func init() {
	APIIntance = helper.NewAPI()

	APIIntance.AddResourceWithWrapper(&KittenResource{}, wrapper, "/api/kittens")
}

func Mux() *http.ServeMux {
	return APIIntance.Mux()
}

func wrapper(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	})
}

type KittenResource struct{}

func (res *KittenResource) Get(r *http.Request) (int, interface{}, http.Header) {
	return 200, map[string]interface{}{"kittens": []map[string]interface{}{
		map[string]interface{}{"id": 1, "name": "Bobby", "picture": "http://placekitten.com/200/200"},
		map[string]interface{}{"id": 2, "name": "Wally", "picture": "http://placekitten.com/200/200"},
	}}, nil
}
