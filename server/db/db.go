package db

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type schemas struct {
	Version  int
	Name     string
	Checksum string
}

// NewDB : Opens a new connection to a database, either to a local test env or by reading a property file.
// 		   This also performs migrations found in sql files
func NewDB(debug bool) *sqlx.DB {
	var (
		conn   string
		folder string
	)
	if debug {
		// new in memory db
		conn = "root:testing@(phoos-db:3306)/phoosball"
		folder = "db/test"
		log.Printf("Using debug connection: %s", conn)
	} else {
		p, err := ioutil.ReadFile("conf/db_conn.txt")
		if err != nil {
			panic(err)
		}
		conn = strings.TrimSuffix(string(p), "\n")
		log.Println(conn)
		folder = "db/prod"
	}

	database := sqlx.MustOpen("mysql", conn)
	performMigrations(database, folder)
	return database
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
