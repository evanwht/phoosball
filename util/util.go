package util

import (
	"database/sql"
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

// DbHandler : wraps a function handler from net/http with a sql.DB parameter
func DbHandler(db *sql.DB, f func(*sql.DB, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(db, w, r)
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