package util

import (
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Env : environment variables that should be shared between routes but created ony once
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
