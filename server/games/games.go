package games

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/evanwht/phoosball/server/players"
	"github.com/evanwht/phoosball/server/util"
	"github.com/gorilla/mux"
)

// Event : event that happens during a game
type Event struct {
	Type int `json:"type"`
	By   int `json:"by"`
	On   int `json:"on"`
}

// Team : Grouping of players that represent one half of a phoosball game
type Team struct {
	Offense  players.Player `json:"offsense"` // player that started on offense
	Deffense players.Player `json:"defense"`  // player that started on defense
}

// Game : data representing a game
type Game struct {
	ID         int     `json:"id"`
	Played     string  `json:"played"`
	Team1      Team    `json:"team1"`
	Team2      Team    `json:"team2"`
	Team1Half  int     `json:"team1Half"`
	Team2Half  int     `json:"team2Half"`
	Team1Final int     `json:"team1Final"`
	Team2Final int     `json:"team2Final"`
	Events     []Event `json:"events"`
}

// GameDB : struct for storing games in memory gotten from the DB
type GameDB struct {
	ID         int    `db:"id"`
	Played     string `db:"played"`
	Team1D     string `db:"team_1_d"`
	Team1DId   int    `db:"team_1_d_id"`
	Team1O     string `db:"team_1_o"`
	Team1OId   int    `db:"team_1_o_id"`
	Team2D     string `db:"team_2_d"`
	Team2DId   int    `db:"team_2_d_id"`
	Team2O     string `db:"team_2_o"`
	Team2OId   int    `db:"team_2_o_id"`
	Team1Half  int    `db:"team_1_half"`
	Team2Half  int    `db:"team_2_half"`
	Team1Final int    `db:"team_1_final"`
	Team2Final int    `db:"team_2_final"`
}

// Endpoint : endpoint that returns JSON information about games
type Endpoint struct {
	Env *util.Env
}

// Get : get all games
func (e Endpoint) Get(r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	var (
		toReturn interface{}
		err      error
	)
	if ok {
		dbGame := GameDB{}
		err = e.Env.DB.Get(&dbGame, "select id, DATE(game_date) as played, team_1_d, team_1_d_id, team_1_o, team_1_o_id, team_2_d, team_2_d_id, team_2_o, team_2_o_id, team_1_half, team_2_half, team_1_final, team_2_final from last_games WHERE id = ?;", id)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}

		toReturn = copyDbGame(&dbGame)
	} else {
		dbGames := []GameDB{}
		err = e.Env.DB.Select(&dbGames, "select id, DATE(game_date) as played, team_1_d, team_1_d_id, team_1_o, team_1_o_id, team_2_d, team_2_d_id, team_2_o, team_2_o_id, team_1_half, team_2_half, team_1_final, team_2_final from last_games;")
		if err != nil {
			return "", http.StatusInternalServerError, err
		}

		games := []Game{}
		for _, dbGame := range dbGames {
			games = append(games, copyDbGame(&dbGame))
		}
		toReturn = games
	}

	return toReturn, http.StatusOK, nil
}

// There has got to be a better way thatn copying to a temp struct
func copyDbGame(dbGame *GameDB) Game {
	return Game{
		ID:     dbGame.ID,
		Played: dbGame.Played,
		Team1: Team{
			Offense: players.Player{
				ID:   dbGame.Team1OId,
				Name: dbGame.Team1O,
			},
			Deffense: players.Player{
				ID:   dbGame.Team1DId,
				Name: dbGame.Team1D,
			},
		},
		Team2: Team{
			Offense: players.Player{
				ID:   dbGame.Team2OId,
				Name: dbGame.Team2O,
			},
			Deffense: players.Player{
				ID:   dbGame.Team2DId,
				Name: dbGame.Team2D,
			},
		},
		Team1Half:  dbGame.Team1Half,
		Team2Half:  dbGame.Team2Half,
		Team1Final: dbGame.Team1Final,
		Team2Final: dbGame.Team2Final,
	}
}

// Post : creates a new game from the json request data
func (e Endpoint) Post(r *http.Request) (interface{}, int, error) {
	var g Game
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&g)
	if err != nil {
		return "", 0, err
	}
	insertSQL := `INSERT INTO games 
			(team_1_p1, team_1_p2, team_2_p1, team_2_p2, team_1_half, team_2_half, team_1_final, team_2_final) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	var res sql.Result
	if g.Team1Final > g.Team2Final {
		res, err = e.Env.DB.Exec(insertSQL, g.Team1.Offense.ID, g.Team1.Deffense.ID, g.Team2.Offense.ID, g.Team2.Deffense.ID, g.Team1Half, g.Team2Half, g.Team1Final, g.Team2Final)
	} else {
		res, err = e.Env.DB.Exec(insertSQL, g.Team2.Offense.ID, g.Team2.Deffense.ID, g.Team1.Offense.ID, g.Team1.Deffense.ID, g.Team2Half, g.Team1Half, g.Team2Final, g.Team1Final)
	}
	if err != nil {
		return "", 0, err
	}
	id, err := res.LastInsertId()
	if err != nil || id <= 0 {
		return "", 0, err
	}
	return strconv.Itoa(int(id)), http.StatusOK, nil
}

// Put : saves a PUT request to alter a games data
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

	var g Game
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&g)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

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
		return "", http.StatusInternalServerError, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return "", http.StatusInternalServerError, err
	}

	return strconv.Itoa(int(rowCnt)), http.StatusOK, nil
}

// Delete : delete a game
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

	insertSQL := `DELETE FROM games WHERE id = ?;`
	res, err := e.Env.DB.Exec(insertSQL, id)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if rows <= 0 {
		return "", http.StatusNotFound, errors.New("game not found")
	}
	return strconv.Itoa(int(rows)), http.StatusOK, nil
}
