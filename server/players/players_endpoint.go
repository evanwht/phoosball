package players

import (
	"net/http"
	"strconv"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"
)

// PlayersEndpoint : endpoint that returns JSON data about many players
type PlayersEndpoint struct {
	Env *util.Env
}

// Get : gets all players
func (e PlayersEndpoint) Get(r *gin.Context) {
	var players []Player
	if err := e.Env.DB.Select(&players, "SELECT id, name, display_name, email FROM players;"); err != nil {
		r.Error(err)
		return
	}
	r.JSON(http.StatusOK, players)
}

// Post : creates a new player
func (e PlayersEndpoint) Post(r *gin.Context) {
	var p Player
	r.BindJSON(&p)
	if len(p.Name) <= 0 {
		r.String(http.StatusBadRequest, "name required")
		return
	}

	insertSQL := `INSERT INTO players (name, display_name, email) VALUES (?, ?, ?);`
	res, err := e.Env.DB.Exec(insertSQL, p.Name, p.NickName, p.Email)
	if err != nil {
		r.Error(err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil || id <= 0 {
		r.Error(err)
		return
	}
	r.String(http.StatusOK, strconv.Itoa(int(id)))
}
