package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type page struct {
	Title     string
	Body      template.HTML
	Players   []string
	GoalTypes []string
}

var templates = template.Must(template.ParseFiles("template.html"))

func setContentType(w http.ResponseWriter, title string) string {
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func loadHTML(fileName string) (*page, error) {
	var (
		body []byte
		err  error
	)
	if templateName := strings.Split(fileName, ".")[0] + "_template.html"; fileExists(templateName) {
		body, err = ioutil.ReadFile(templateName)
	} else if fileExists(fileName) {
		body, err = ioutil.ReadFile(fileName)
	}
	if err != nil {
		return nil, err
	}
	return &page{Title: fileName, Body: template.HTML(body)}, nil
}

func serveTemplate(w http.ResponseWriter, p *page) {
	err := templates.ExecuteTemplate(w, "template.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dbHandler(db *sql.DB, f func(*sql.DB, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(db, w, r)
	}
}

func createGameInfo(db *sql.DB) *gameInfo {
	var (
		id          int
		displayName string
		name        string
		names       []string
		events      []string
	)
	// get user info
	rows, err := db.Query("select id, name, display_name from players;")
	if err != nil {
		// do nothing
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &displayName)
			if err != nil {
				displayName = ""
			}
			names = append(names, "\""+displayName+" ("+name+")\",")
		}
		names = append(names, "\"New Player\"")
	}

	// get event type info
	rows, err = db.Query("select * from event_types;")
	if err != nil {
		// do nothing
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err == nil {
				events = append(events, "\""+name+"\",")
			} else {
				// load page with error message
			}
		}
		events = append(events, "\"New Type\"")
	}
	return &gameInfo{Players: names, GoalTypes: events}
}

type gameInfo struct {
	Players   []string
	GoalTypes []string
}

func renderGamePage(db *sql.DB, w http.ResponseWriter, r *http.Request) template.HTML {
	var g = createGameInfo(db)

	t, err := template.ParseFiles("game_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buff.String())
}

func defaultHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	if len(title) == 0 {
		title = "index"
	}
	switch fileType := setContentType(w, title); fileType {
	case "text/html", "text/plain":
		p, err := loadHTML(title + "_template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if title == "game" {
				p.Body = renderGamePage(db, w, r)
			}
			serveTemplate(w, p)
		}
	default:
		p, err := ioutil.ReadFile(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "%s", p)
		}
	}
}

func addNewPlayer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.PostForm) > 0 {
		p, err := template.ParseFiles("account.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var buff bytes.Buffer
		if err = p.Execute(&buff, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		serveTemplate(w, &page{Title: "Account", Body: template.HTML(buff.String())})
	}
}

func main() {
	p, err := ioutil.ReadFile("../db_conn.txt")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", strings.TrimSuffix(string(p), "\n"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", dbHandler(db, defaultHandler))
	http.HandleFunc("/save_player", dbHandler(db, addNewPlayer))
	log.Fatal(http.ListenAndServe(":3032", nil))
}
