// Runs a REST server for a phoosball website. Includes a react SPA and RESTful JSON api server
package main

import (
	"flag"
	"net/http"

	"github.com/evanwht/phoosball/server/db"
	"github.com/evanwht/phoosball/server/games"
	"github.com/evanwht/phoosball/server/players"
	"github.com/evanwht/phoosball/server/standings"
	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	// get command line argument for if this is in test mode
	boolPtr := flag.Bool("debug", false, "inits a sqlite db in memory for local testing")
	flag.Parse()

	database := db.NewDB(*boolPtr)
	defer database.Close()
	env := &util.Env{DB: database}

	r := gin.Default()

	gamesEndpoint := games.GamesEndpoint{Env: env}
	gameEndpoint := games.GameEndpoint{Env: env}
	playersEndpoint := players.PlayersEndpoint{Env: env}
	playerEndpoint := players.PlayerEndpoint{Env: env}
	standingsEndpoint := standings.StandingsEndpoint{Env: env}
	standingEndpoint := standings.StandingEndpoint{Env: env}
	api := r.Group("/api")
	{
		games := api.Group("/games")
		{
			games.GET("", gamesEndpoint.Get)
			games.POST("", gamesEndpoint.Post)
			games.GET("/id:[0-9]+", gameEndpoint.Get)
			games.PUT("/id:[0-9]+", gameEndpoint.Put)
			games.DELETE("/id:[0-9]+", gameEndpoint.Delete)
		}
		players := api.Group("/players")
		{
			players.GET("", playersEndpoint.Get)
			players.POST("", playersEndpoint.Post)
			players.GET("/id:[0-9]+", playerEndpoint.Get)
			players.PUT("/id:[0-9]+", playerEndpoint.Put)
			players.DELETE("/id:[0-9]+", playerEndpoint.Delete)
		}
		standings := api.Group("/standings")
		{
			standings.GET("", standingsEndpoint.Get)
			standings.GET("/id:[0-9]+", standingEndpoint.Get)
		}
	}

	r.Static("static", ".")

	r.Run()
}
