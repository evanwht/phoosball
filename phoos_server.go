package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/evanwht1/phoosball/gopages"
	"github.com/evanwht1/phoosball/util"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
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

type schemas struct {
	Version  int
	Name     string
	Checksum string
}

func performMigrations(db *sqlx.DB, folder string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Slice(sqlFiles, func(i, j int) bool {
		ind, _ := strconv.Atoi(sqlFiles[i][0:1])
		jnd, _ := strconv.Atoi(sqlFiles[j][0:1])
		return ind < jnd
	})
	var filesApplied = []schemas{}
	err = db.Select(&filesApplied, "select version, name, checksum from schema_history;")
	if err != nil {
		panic(err)
	}
	for _, file := range sqlFiles {
		query, err := ioutil.ReadFile(folder + "/" + file)
		if err != nil {
			panic(err)
		}
		parts := strings.Split(file, ".")
		v, _ := strconv.Atoi(parts[0])
		sum := md5.Sum(query)
		sf := hex.EncodeToString(sum[:])
		if !(contains(filesApplied, v, file, sf)) {
			fmt.Printf("applying %s\n", file)
			inserts := string(query)
			tx, err := db.Begin()
			if err != nil {
				panic(err)
			}
			stmts := strings.Split(inserts, ";\n")
			for _, in := range stmts {
				if _, err := tx.Exec(in); err != nil {
					panic(err)
				}
			}
			
			_, err = tx.Exec("INSERT INTO schema_history (version, description, name, checksum) VALUES (?, ?, ?, ?);", v, parts[1], file, sf)
			if err != nil {
				log.Print("rolling back changes")
				tx.Rollback()
				log.Fatal(err)
			} else {
				log.Print("committing changes")
				tx.Commit()
			}
		}
	}
}

// contains : checks if a string is contained in an array of string
func contains(arr []schemas, ver int, name string, checksum string) bool {
	for _, s := range arr {
		if s.Version == ver {
			if s.Name != name || s.Checksum != checksum {
				log.Printf("Schema file changed from what was applied to db: %s", name)
				panic(errors.New("Changes to already applied DB Migration files detected"))
			}
			return true
		}
	}
	return false
}

func main() {
	// get command line argument for if this is in test mode
	boolPtr := flag.Bool("debug", false, "inits a sqlite db in memory for local testing")
	flag.Parse()
	var (
		env    *util.Env
		conn   string
		folder string
	)
	if *boolPtr {
		// new in memory db
		conn = "root:testing@(127.0.0.1:3306)/phoosball"
		folder = "db/test"
		log.Printf("Using debug connection: %s", conn)
	} else {
		p, err := ioutil.ReadFile("db_conn.txt")
		if err != nil {
			panic(err)
		}
		conn = strings.TrimSuffix(string(p), "\n")
		folder = "db/prod"
	}

	db := sqlx.MustOpen("mysql", conn)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	defer db.Close()
	performMigrations(db, folder)
	env = &util.Env{DB: db}

	r := mux.NewRouter()

	r.HandleFunc("/players", util.Chain(env, playersHandler, util.Methods("GET"), util.Headers()))

	player := r.PathPrefix("/player").Subrouter()
	player.PathPrefix("/edit").Handler(util.Chain(env, playerHandler, util.Methods("PUT"), util.Headers()))
	player.Handle("", util.Chain(env, playerHandler, util.Methods("GET", "POST"), util.Headers()))

	r.HandleFunc("/games", util.Chain(env, gamesHandler, util.Headers())).Methods("GET")

	game := r.PathPrefix("/game").Subrouter()
	game.Handle("", util.Chain(env, gameHandler, util.Methods("GET", "POST"), util.Headers()))
	game.PathPrefix("/edit").Handler(util.Chain(env, gopages.SaveGameEdit, util.Methods("PUT"), util.Headers()))

	r.HandleFunc("/favicon.ico", faviconHandler)

	// requests for other paths with incorrect methods will get here and try to render something instead
	// of failing in their correct handler. Short coming of Gorilla IMO
	r.PathPrefix("/").Handler(util.Chain(env, defaultHandler, util.Headers()))

	log.Fatal(http.ListenAndServe(":3032", r))
}
