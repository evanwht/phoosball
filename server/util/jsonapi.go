// Package util provides helper methods for creating a JSON http.Handler wrapper for a JSON REST endpoint
package util

import (
	"encoding/json"
	"net/http"
)

// JSONAPIEndpoint : RESTful endpoit that returns JSON data for each method
type JSONAPIEndpoint interface {
	Get(r *http.Request) (interface{}, int, error)
	Post(r *http.Request) (interface{}, int, error)
	Put(r *http.Request) (interface{}, int, error)
	Delete(r *http.Request) (interface{}, int, error)
}

// JSONHandler : wrapper of http.Handler for a JSONEndpoint
type JSONHandler struct {
	Endpoint JSONAPIEndpoint
}

func (j JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var (
		resp      interface{}
		errStatus int
		err       error
	)
	switch r.Method {
	case "GET":
		resp, errStatus, err = j.Endpoint.Get(r)
	case "PUT":
		resp, errStatus, err = j.Endpoint.Put(r)
	case "POST":
		resp, errStatus, err = j.Endpoint.Post(r)
	case "DELETE":
		resp, errStatus, err = j.Endpoint.Delete(r)
	default:
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
	}
	if err != nil {
		http.Error(w, err.Error(), errStatus)
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}
