package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/evanwht1/phoosball/gopages"
	"github.com/evanwht1/phoosball/util"

	_ "github.com/go-sql-driver/mysql"
)

var templates = template.Must(template.ParseFiles("webpage/template.html"))

func loadHTML(fileName string) (*util.Page, error) {
	var (
		body []byte
		err  error
	)
	if templateName := strings.Split(fileName, ".")[0] + "_template.html"; util.FileExists(templateName) {
		body, err = ioutil.ReadFile(templateName)
	} else if util.FileExists(fileName) {
		body, err = ioutil.ReadFile(fileName)
	}
	if err != nil {
		return nil, err
	}
	return &util.Page{Title: fileName, Body: template.HTML(body)}, nil
}

func serveTemplate(w http.ResponseWriter, p *util.Page) {
	err := templates.ExecuteTemplate(w, "template.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func defaultHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	title := "webpage" + r.URL.Path
	if len(title) == 0 {
		title = "index"
	}
	switch fileType := util.SetContentType(w, title); fileType {
	case "text/html", "text/plain":
		p, err := loadHTML(title + "_template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if title == "webpage/game" {
				p.Body, err = gopages.RenderGamePage(db, w, r)
				if err != nil {
					return
				}
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

// func defaultHandler(f func(*sql.DB, http.ResponseWriter, *http.Request)) {
// 	serveTemplate(w, f())
// }

func main() {
	p, err := ioutil.ReadFile("db_conn.txt")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", strings.TrimSuffix(string(p), "\n"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", util.DbHandler(db, defaultHandler))
	// http.HandleFunc("/save_player", util.DbHandler(db, booty(gopages.AddNewPlayer)))
	log.Fatal(http.ListenAndServe(":3032", nil))
}
