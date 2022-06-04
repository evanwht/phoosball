package standings

import (
	"net/http"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gin-gonic/gin"
)

// Endpoint : endpoint that returns JSON data for player standings
type StandingEndpoint struct {
	Env *util.Env
}

// Get : get standings of all players
func (e StandingEndpoint) Get(r *gin.Context) {
	id := r.Param("id")
	var standing Standing
	if err := e.Env.DB.Get(&standing, "select * from overall_standings WHERE id = ?;", id); err != nil {
		r.Error(err)
		return
	}
	r.JSON(http.StatusOK, standing)
}
