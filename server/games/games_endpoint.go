package games

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"
)

// Endpoint : endpoint that returns JSON information about games
type GamesEndpoint struct {
	Env *util.Env
}

// Get : get all games
func (e GamesEndpoint) Get(r *gin.Context) {
	var dbGames []GameDB
	if err := e.Env.DB.Select(&dbGames, "select id, DATE(game_date) as played, team_1_d, team_1_d_id, team_1_o, team_1_o_id, team_2_d, team_2_d_id, team_2_o, team_2_o_id, team_1_half, team_2_half, team_1_final, team_2_final from last_games;"); err != nil {
		r.Error(err)
		return
	}

	games := []Game{}
	for _, dbGame := range dbGames {
		games = append(games, copyDbGame(&dbGame))
	}
	r.JSON(http.StatusOK, games)
}

// Post : creates a new game from the json request data
func (e GamesEndpoint) Post(r *gin.Context) {
	var g Game
	r.BindJSON(&g)

	insertSQL := `INSERT INTO games 
			(team_1_p1, team_1_p2, team_2_p1, team_2_p2, team_1_half, team_2_half, team_1_final, team_2_final) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	var (
		res sql.Result
		err error
	)
	if g.Team1Final > g.Team2Final {
		res, err = e.Env.DB.Exec(insertSQL, g.Team1.Offense.ID, g.Team1.Deffense.ID, g.Team2.Offense.ID, g.Team2.Deffense.ID, g.Team1Half, g.Team2Half, g.Team1Final, g.Team2Final)
	} else {
		res, err = e.Env.DB.Exec(insertSQL, g.Team2.Offense.ID, g.Team2.Deffense.ID, g.Team1.Offense.ID, g.Team1.Deffense.ID, g.Team2Half, g.Team1Half, g.Team2Final, g.Team1Final)
	}
	if err != nil {
		r.Error(err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil || id <= 0 {
		r.Error(err)
		return
	}
	r.JSON(http.StatusOK, strconv.Itoa(int(id)))
}
