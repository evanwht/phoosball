package gopages

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/evanwht1/phoosball/util"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var gameRowTemplate = template.Must(template.ParseFiles("webpage/games_view/game_row.html"))

func buildRow(gd gameDBData) (string, error) {
	var data gameData
	if gd.T1final > gd.T2final {
		data = gameData{gd, "success", "danger"}
	} else {
		data = gameData{gd, "danger", "success"}
	}
	var buff bytes.Buffer
	if err := gameRowTemplate.Execute(&buff, data); err != nil {
		return "", err
	}
	return buff.String(), nil
}

func getGames(db *sqlx.DB) *gamesInfo {
	var tableRows []string
	games := []gameDBData{}
	err := db.Select(&games, "select id, DATE(game_date) as game_date, team_1_Defense, team_1_Offense, team_2_Defense, team_2_Offense, team_1_half, team_2_half, team_1_final, team_2_final from last_games;")
	if err != nil {
		log.Fatal(err)
	} else if len(games[0].Date) == 0 {
		log.Fatal("Select returned nothing")
	} else {
		for _, game := range games {
			st, err := buildRow(game)
			if err != nil {
				log.Fatal(err)
			}
			tableRows = append(tableRows, st)
		}
	}
	return &gamesInfo{Games: template.HTML(strings.Join(tableRows, "\n"))}
}

type gamesInfo struct {
	Games template.HTML
}

// RenderGamesPage : gets data from db to show last 10 games played
func RenderGamesPage(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	t, err := template.ParseFiles("webpage/games_view/games_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	g := getGames(db)
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	return template.HTML(buff.String()), nil
}

// SaveGameEdit : saves a PUT request to alter a games data
func SaveGameEdit(env *util.Env, w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
		var t gameData
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			tx, err := env.DB.Begin()
			var fail bool
			if err != nil {
				log.Fatal(err)
				fail = true
			} else {
				stmt, err := tx.Prepare(`UPDATE games SET
									team_1_p1 = ?, team_1_p2 = ?, team_2_p1 = ?, team_2_p2 = ?,
									team_1_half = ?, team_2_half = ?, team_1_final = ?, team_2_final = ? 
									WHERE id = ?;`)
				if err != nil {
					log.Println(err)
					fail = true
				} else {
					var res sql.Result
					if t.T1final > t.T2final {
						res, err = stmt.Exec(t.T1pd, t.T1po, t.T2pd, t.T2po, t.T1half, t.T2half, t.T1final, t.T2final, t.ID)
					} else {
						res, err = stmt.Exec(t.T2pd, t.T2po, t.T1pd, t.T1po, t.T2half, t.T1half, t.T2final, t.T1final, t.ID)
					}
					if err != nil {
						log.Println(err)
						fail = true
					}
					rowCnt, err := res.RowsAffected()
					if err != nil || rowCnt != 1 {
						fail = true
					}
				}
				if fail {
					tx.Rollback()
					http.Error(w, "Error", http.StatusInternalServerError)
				} else {
					tx.Commit()
					fmt.Fprint(w, "Saved")
				}
				stmt.Close()
			}
		}
	} else {
		http.Error(w, "NOT ALLOWED", http.StatusMethodNotAllowed)
	}
}
