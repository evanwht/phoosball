package gopages

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var gameRowTemplate = template.Must(template.ParseFiles("webpage/games_view/game_row.html"))

func buildRow(gd gameData) (string, error) {
	if gd.T1final > gd.T2final {
		gd.Team1Class = "success"
		gd.Team2Class = "danger"
	} else {
		gd.Team1Class = "danger"
		gd.Team2Class = "success"
	}
	var buff bytes.Buffer
	if err := gameRowTemplate.Execute(&buff, gd); err != nil {
		return "", err
	}
	return buff.String(), nil
}

func getGames(db *sql.DB) *gamesInfo {
	var tableRows []string
	rows, err := db.Query("select * from last_games;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var game gameData
			err := rows.Scan(&game.Date, &game.T1pd, &game.T1po, &game.T2pd, &game.T2po, &game.T1half, &game.T2half, &game.T1final, &game.T2final)
			game.Date = game.Date[:11]
			if err != nil {
				log.Fatal(err)
			} else {
				st, err := buildRow(game)
				if err != nil {
					log.Fatal(err)
				}
				tableRows = append(tableRows, st)
			}
		}
		rows.Close()
	}
	return &gamesInfo{Games: template.HTML(strings.Join(tableRows, "\n"))}
}

type gamesInfo struct {
	Games template.HTML
}

// RenderGamesPage : gets data from db to show last 10 games played
func RenderGamesPage(db *sql.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	t, err := template.ParseFiles("webpage/games_template.html")
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
func SaveGameEdit(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var t gameData
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			// TODO create db update statement
			log.Print(t)
		}
	} else {
		http.Error(w, "NOT ALLOWED", http.StatusMethodNotAllowed)
	}
}
