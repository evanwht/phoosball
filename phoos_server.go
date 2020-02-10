package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/evanwht1/phoosball/gopages"
	"github.com/evanwht1/phoosball/util"
	"github.com/gorilla/mux"

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

func gameHandler(env *util.Env, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	p, err := loadHTML("webpage/game_input/game_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		p.Body, err = gopages.RenderGamePage(env.DB, w, r)
		if err != nil {
			return
		}
		serveTemplate(w, p)
	}
}

func gamesHandler(env *util.Env, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	body, err := gopages.RenderGamesPage(env.DB, w, r)
	if err != nil {
		return
	}
	serveTemplate(w, &util.Page{Title: "Games", Body: body})
}

func playerHandler(env *util.Env, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	body, err := gopages.RenderPlayerPage(env.DB, w, r)
	if err != nil {
		return
	}
	serveTemplate(w, &util.Page{Title: "Player", Body: body})
}

func playersHandler(env *util.Env, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	body := gopages.GetAllPlayers(env.DB)
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(b))
}

func indexHandler(env *util.Env, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	body, err := gopages.RenderStandingsPage(env.DB, w, r)
	if err != nil {
		return
	}
	serveTemplate(w, &util.Page{Title: "Standings", Body: body})
}

func defaultHandler(env *util.Env, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		indexHandler(env, w, r)
		return
	}
	title := "webpage" + r.URL.Path
	switch fileType := util.SetContentType(w, title); fileType {
	case "text/html", "text/plain":
		p, err := loadHTML(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			serveTemplate(w, p)
		}
		break
	default:
		p, err := ioutil.ReadFile(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "%s", p)
		}
		break
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
	env := &util.Env{DB: db}

	r := mux.NewRouter()
	
	r.HandleFunc("/players", util.Chain(env, playersHandler, util.Headers())).Methods("GET")
	
	player := r.PathPrefix("/player").Subrouter()
	player.PathPrefix("/edit").Methods("PUT").Handler(util.Chain(env, playerHandler, util.Headers()))
	player.Methods("GET", "POST").Handler(util.Chain(env, playerHandler, util.Headers()))

	
	r.HandleFunc("/games", util.Chain(env, gamesHandler, util.Headers())).Methods("GET")

	game := r.PathPrefix("/game").Subrouter()
	game.Methods("GET", "POST").Handler(util.Chain(env, gameHandler, util.Headers()))
	game.PathPrefix("/edit").Methods("PUT").Handler(util.Chain(env, gopages.SaveGameEdit, util.Headers()))


	r.HandleFunc("/favicon.ico", faviconHandler)

	r.PathPrefix("/").Handler(util.Chain(env, defaultHandler, util.Headers()))

	log.Fatal(http.ListenAndServe(":3032", r))
}
