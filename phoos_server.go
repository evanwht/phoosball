package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"encoding/json"

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
	templateName := strings.Split(fileName, ".")[0] + "_template.html"
	if util.FileExists(templateName) {
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

func gameHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	w.Header().Set("Content-Language", "en-US")
	p, err := loadHTML("webpage/game_input/game_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		p.Body, err = gopages.RenderGamePage(db, w, r)
		if err != nil {
			return
		}
		serveTemplate(w, p)
	}
}

func gamesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	w.Header().Set("Content-Language", "en-US")
	body, err := gopages.RenderGamesPage(db, w, r)
	if err != nil {
		return
	}
	serveTemplate(w, &util.Page{Title: "Games", Body: body})
}

func playerHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	w.Header().Set("Content-Language", "en-US")
	body, err := gopages.RenderPlayerPage(db, w, r)
	if err != nil {
		return
	}
	serveTemplate(w, &util.Page{Title: "Player", Body: body})
}

func playersHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	w.Header().Set("Content-Language", "en-US")
	body := gopages.GetAllPlayers(db)
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(b))
}

func indexHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	w.Header().Set("Content-Language", "en-US")
	body, err := gopages.RenderStandingsPage(db, w, r)
	if err != nil {
		return
	}
	serveTemplate(w, &util.Page{Title: "Standings", Body: body})
}

func defaultHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		indexHandler(db, w, r)
		return
	}
	title := "webpage" + r.URL.Path
	w.Header().Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	w.Header().Set("Content-Language", "en-US")
	switch fileType := util.SetContentType(w, title); fileType {
	case "text/html", "text/plain":
		p, err := loadHTML(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "webpage/favicon/favicon.ico")
}

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
	http.HandleFunc("/player", util.DbHandler(db, playerHandler))
	http.HandleFunc("/players", util.DbHandler(db, playersHandler))
	http.HandleFunc("/game", util.DbHandler(db, gameHandler))
	http.HandleFunc("/games", util.DbHandler(db, gamesHandler))
	http.HandleFunc("/edit_game", util.DbHandler(db, gopages.SaveGameEdit))
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(":3032", nil))
}
