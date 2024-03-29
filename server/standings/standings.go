package standings

import (
	"net/http"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"
)

// Endpoint : endpoint that returns JSON data for player standings
type StandingsEndpoint struct {
	Env *util.Env
}

// Get : get standings of all players
func (e StandingsEndpoint) Get(r *gin.Context) {
	var standings []Standing
	if err := e.Env.DB.Select(&standings, "select * from overall_standings;"); err != nil {
		r.Error(err)
		return
	}
	r.JSON(http.StatusOK, standings)
}
