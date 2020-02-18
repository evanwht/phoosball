package util

import (
	"github.com/jmoiron/sqlx"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// Page : holder of data for pages
type Page struct {
	Title     string
	Body      template.HTML
	Players   []string
	GoalTypes []string
}

// Env : envorionment variables that should be shared between routes but created ony once
type Env struct {
	DB *sqlx.DB
}

// Middleware : does this to request before route function
type Middleware func(RouteFunc) RouteFunc

// RouteFunc :
type RouteFunc func(*Env, http.ResponseWriter, *http.Request)

// Headers : writes common headers to all routes
func Headers() Middleware {
	return func(rf RouteFunc) RouteFunc {
		return func(env *Env, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
			w.Header().Set("Content-Language", "en")
			rf(env, w, r)
		}
	}
}

// Methods : writes common headers to all routes
func Methods(methods ...string) Middleware {
	return func(rf RouteFunc) RouteFunc {
		return func(env *Env, w http.ResponseWriter, r *http.Request) {
			// This is one example of how go is idiotic. "It is trivial to write your own contains method"
			// If it is trivial, why doesn't the language just provide it since you don't give support
			// for generics, thus making EVERY SINGLE PERSON write the same lines of code to see if an
			// array of {type} contains a specific value. Writers of go are egotistical ass holes
			for _, m := range methods {
				if r.Method == m {
					rf(env, w, r)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

// DBRoute : wraps a function handler from net/http with Env parameters
func DBRoute(env *Env, rf RouteFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rf(env, w, r)
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(env *Env, f RouteFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		f(env, w, r)
	}
}

// SetContentType : sets the content type header based on the file extension
func SetContentType(w http.ResponseWriter, title string) string {
	fileType := append(strings.Split(title, "."), "")[1]
	switch fileType {
	case "css":
		w.Header().Set("content-type", "text/css")
		return "text/css"
	case "svg":
		w.Header().Set("content-type", "image/svg+xml")
		return "image/svg+xml"
	case "png":
		w.Header().Set("content-type", "image/png")
		return "image/png"
	case "jpg":
		w.Header().Set("content-type", "image/jpg")
		return "image/jpg"
	case "js":
		w.Header().Set("content-type", "text/javascript")
		return "text/javascript"
	default:
		w.Header().Set("content-type", "text/html")
		return "text/html"
	}
}

// FileExists : checks if a file exits or not
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// HTMLOption : wraps the id and text in an html option tag
func HTMLOption(id string, text string) string {
	return "<option value=\"" + id + "\">" + text + "</opton>"
}

// Contains : checks if a string is contained in an array of string
func Contains(arr []string, str string) bool {
	for _, s := range arr {
		if (s == str) {
			return true
		}
	}
	return false;
}