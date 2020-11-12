// Runs a REST server for a phoosball website. Includes a react SPA and RESTful JSON api server
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/evanwht/phoosball/server/db"
	"github.com/evanwht/phoosball/server/games"
	"github.com/evanwht/phoosball/server/players"
	"github.com/evanwht/phoosball/server/standings"
	"github.com/evanwht/phoosball/server/util"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	// get command line argument for if this is in test mode
	boolPtr := flag.Bool("debug", false, "inits a sqlite db in memory for local testing")
	flag.Parse()

	database := db.NewDB(*boolPtr)
	defer database.Close()
	env := &util.Env{DB: database}

	r := mux.NewRouter()

	gamesHandler := util.JSONHandler{
		Endpoint: games.Endpoint{Env: env},
	}
	r.PathPrefix("/api/games/{id:[0-9]+}").
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete).
		Handler(gamesHandler)
	r.PathPrefix("/api/games").
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete).
		Handler(gamesHandler)

	playersHandler := util.JSONHandler{
		Endpoint: players.Endpoint{Env: env},
	}
	r.PathPrefix("/api/players/{id:[0-9]+}").
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete).
		Handler(playersHandler)
	r.PathPrefix("/api/players").
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete).
		Handler(playersHandler)

	standingsHandler := util.JSONHandler{
		Endpoint: standings.Endpoint{Env: env},
	}
	r.PathPrefix("/api/standings/{id:[0-9]+}").
		Methods(http.MethodGet).
		Handler(standingsHandler)
	r.PathPrefix("/api/standings").
		Methods(http.MethodGet).
		Handler(standingsHandler)

	// staticFileHandler := http.StripPrefix("/static/", http.FileServer((http.Dir("/static"))))
	// r.PathPrefix("/static/").Handler(staticFileHandler)

	// always put index router last as the first route the request matches will execute that router
	// fileHandler := http.FileServer(http.Dir("/public"))
	// r.PathPrefix("/").Handler(fileHandler)

	spa := spaHandler{staticPath: "public", indexPath: "index.html"}

	r.PathPrefix("/").Handler(spa)

	r.Use(loggingMiddleware)

	srv := &http.Server{
		Addr: ":3032",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		log.Printf("Starting server %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
