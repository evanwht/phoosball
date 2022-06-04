package games

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"
)

// Endpoint : endpoint that returns JSON information about games
type GameEndpoint struct {
	Env *util.Env
}

// Get : get all games
func (e GameEndpoint) Get(r *gin.Context) {
	id := r.Param("id")
	var (
		toReturn interface{}
		err      error
	)
	dbGame := GameDB{}
	err = e.Env.DB.Get(&dbGame, "select id, DATE(game_date) as played, team_1_d, team_1_d_id, team_1_o, team_1_o_id, team_2_d, team_2_d_id, team_2_o, team_2_o_id, team_1_half, team_2_half, team_1_final, team_2_final from last_games WHERE id = ?;", id)
	if err != nil {
		r.Error(err)
		return
	}

	toReturn = copyDbGame(&dbGame)
	r.JSON(http.StatusOK, toReturn)
}

// Put : saves a PUT request to alter a games data
func (e GameEndpoint) Put(r *gin.Context) {
	idStr := r.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		r.Error(err)
		return
	}

	var g Game
	r.BindJSON(&g)

	updateGame := `UPDATE games SET
							team_1_p1 = ?, team_1_p2 = ?, team_2_p1 = ?, team_2_p2 = ?,
							team_1_half = ?, team_2_half = ?, team_1_final = ?, team_2_final = ? 
							WHERE id = ?;`
	var res sql.Result
	if g.Team1Final > g.Team2Final {
		res, err = e.Env.DB.Exec(updateGame, g.Team1.Offense.ID, g.Team1.Deffense.ID, g.Team2.Offense.ID, g.Team2.Deffense.ID, g.Team1Half, g.Team2Half, g.Team1Final, g.Team2Final, id)
	} else {
		res, err = e.Env.DB.Exec(updateGame, g.Team2.Offense.ID, g.Team2.Deffense.ID, g.Team1.Offense.ID, g.Team1.Deffense.ID, g.Team2Half, g.Team1Half, g.Team2Final, g.Team1Final, id)
	}

	if err != nil {
		r.Error(err)
		return
	}

	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		r.Error(err)
		return
	}

	r.JSON(http.StatusOK, strconv.Itoa(int(rowCnt)))
}

// Delete : delete a game
func (e GameEndpoint) Delete(r *gin.Context) {
	idStr := r.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		r.Error(err)
		return
	}

	insertSQL := `DELETE FROM games WHERE id = ?;`
	res, err := e.Env.DB.Exec(insertSQL, id)
	if err != nil {
		r.Error(err)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.Error(err)
		return
	}
	if rows <= 0 {
		r.Error(err)
		return
	}
	r.JSON(http.StatusOK, strconv.Itoa(int(rows)))
}
