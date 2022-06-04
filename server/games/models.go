package games

import "github.com/evanwht/phoosball/server/players"

// Event : event that happens during a game
type Event struct {
	Type int `json:"type"`
	By   int `json:"by"`
	On   int `json:"on"`
}

// Team : Grouping of players that represent one half of a phoosball game
type Team struct {
	Offense  players.Player `json:"offense"` // player that started on offense
	Deffense players.Player `json:"defense"` // player that started on defense
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
