package players

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/evanwht/phoosball/server/util"
	"github.com/gorilla/mux"
)

// Player from the db
type Player struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	NickName string `json:"nickname,omitempty" db:"display_name"`
	Email    string `json:"email,omitempty" db:"email"`
}

// Endpoint : endpoint that returns JSON data about players
type Endpoint struct {
	Env *util.Env
}

// Get : gets all players
func (e Endpoint) Get(r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	var (
		toReturn interface{}
		err      error
	)
	if ok {
		player := Player{}
		err = e.Env.DB.Get(&player, "SELECT id, name, display_name, email FROM players WHERE id = ?;", id)
		toReturn = player
	} else {
		players := []Player{}
		err = e.Env.DB.Select(&players, "SELECT id, name, display_name, email FROM players;")
		toReturn = players
	}
	if err != nil {
		return "", http.StatusInternalServerError, nil
	}
	return toReturn, http.StatusOK, nil
}

// Post : creates a new player
func (e Endpoint) Post(r *http.Request) (interface{}, int, error) {
	var p Player
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	if len(p.Name) <= 0 {
		return "", http.StatusBadRequest, errors.New("player name required")
	}

	insertSQL := `INSERT INTO players (name, display_name, email) VALUES (?, ?, ?);`
	res, err := e.Env.DB.Exec(insertSQL, p.Name, p.NickName, p.Email)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	id, err := res.LastInsertId()
	if err != nil || id <= 0 {
		return "", http.StatusInternalServerError, err
	}
	return strconv.Itoa(int(id)), http.StatusOK, nil
}

// Delete : delete a player
func (e Endpoint) Delete(r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	pathID, ok := vars["id"]
	if !ok {
		return "", http.StatusBadRequest, errors.New("Missing id")
	}
	id, err := strconv.Atoi(pathID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	insertSQL := `DELETE FROM players WHERE id = ?;`
	res, err := e.Env.DB.Exec(insertSQL, id)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if rows <= 0 {
		return "", http.StatusNotFound, errors.New("player not found")
	}
	return strconv.Itoa(int(rows)), http.StatusOK, nil
}

// Put : update a players info
func (e Endpoint) Put(r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	pathID, ok := vars["id"]
	if !ok {
		return "", http.StatusBadRequest, errors.New("Missing id")
	}
	id, err := strconv.Atoi(pathID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	var p Player
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&p)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	updatePlayer := `UPDATE players SET name = ?, display_name = ?, email = ? WHERE id = ?;`
	res, err := e.Env.DB.Exec(updatePlayer, p.Name, p.NickName, p.Email, id)

	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return "", http.StatusInternalServerError, err
	}

	return strconv.Itoa(int(rowCnt)), http.StatusOK, nil
}
