package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

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
		// env    *util.Env
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
	// env = &util.Env{DB: db}

	r := mux.NewRouter()

	fileHandler := http.FileServer(http.Dir("../ui/build/public"))
	r.PathPrefix("/").Handler(fileHandler)

	staticFileHandler := http.StripPrefix("/static/", http.FileServer((http.Dir("../ui/build/public/static"))))
	r.PathPrefix("/static/").Handler(staticFileHandler)

	log.Fatal(http.ListenAndServe(":3032", r))
}
