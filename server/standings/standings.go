package standings

import (
	"errors"
	"net/http"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gorilla/mux"
)

// Standing : wins and losses of a player in a context
type Standing struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}

// Endpoint : endpoint that returns JSON data for player standings
type Endpoint struct {
	Env *util.Env
}

// Get : get standings of all players
func (e Endpoint) Get(r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	var (
		toReturn interface{}
		err      error
	)
	if ok {
		standing := Standing{}
		err = e.Env.DB.Get(&standing, "select * from overall_standings WHERE id = ?;", id)
		toReturn = standing
	} else {
		standings := []Standing{}
		err = e.Env.DB.Select(&standings, "select * from overall_standings;")
		toReturn = standings
	}

	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return toReturn, http.StatusOK, nil
}

// Post : N/A
func (e Endpoint) Post(r *http.Request) (interface{}, int, error) {
	return "", http.StatusMethodNotAllowed, errors.New("method not allowed")
}

// Put : N/A
func (e Endpoint) Put(r *http.Request) (interface{}, int, error) {
	return "", http.StatusMethodNotAllowed, errors.New("method not allowed")
}

// Delete : N/A
func (e Endpoint) Delete(r *http.Request) (interface{}, int, error) {
	return "", http.StatusMethodNotAllowed, errors.New("method not allowed")
}
