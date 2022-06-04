package players

import (
	"net/http"
	"strconv"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"
)

// PlayerEndpoint : endpoint that returns JSON data about a single player
type PlayerEndpoint struct {
	Env *util.Env
}

// Get : gets a single player
func (e PlayerEndpoint) Get(r *gin.Context) {
	id := r.Param("id")

	var player Player
	if err := e.Env.DB.Get(&player, "SELECT id, name, display_name, email FROM players WHERE id = ?;", id); err != nil {
		r.Error(err)
		return
	}
	r.JSON(http.StatusOK, player)
}

// Delete : delete a single player
func (e PlayerEndpoint) Delete(r *gin.Context) {
	pathID := r.Param("id")
	id, err := strconv.Atoi(pathID)
	if err != nil {
		r.Error(err)
		return
	}

	insertSQL := `DELETE FROM players WHERE id = ?;`
	res, err := e.Env.DB.Exec(insertSQL, id)
	if err != nil {
		r.Error(err)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil || rows <= 0 {
		r.Error(err)
		return
	}
	r.String(http.StatusOK, strconv.Itoa(int(rows)))
}

// Put : update a single player
func (e PlayerEndpoint) Put(r *gin.Context) {
	pathID := r.Param("id")
	id, err := strconv.Atoi(pathID)
	if err != nil {
		r.Error(err)
		return
	}

	var p Player
	r.BindJSON(&p)

	updatePlayer := `UPDATE players SET name = ?, display_name = ?, email = ? WHERE id = ?;`
	res, err := e.Env.DB.Exec(updatePlayer, p.Name, p.NickName, p.Email, id)

	if err != nil {
		r.Error(err)
		return
	}

	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		r.Error(err)
		return
	}

	r.String(http.StatusOK, strconv.Itoa(int(rowCnt)))
}
