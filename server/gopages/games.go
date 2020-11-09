package gopages

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/evanwht/phoosball/server/server/util"
)

// var gameRowTemplate = template.Must(template.ParseFiles("webpage/games_view/game_row.html"))

// func buildRow(gd gameDBData) (string, error) {
// 	var data gameData
// 	if gd.T1final > gd.T2final {
// 		data = gameData{gd, "success", "danger"}
// 	} else {
// 		data = gameData{gd, "danger", "success"}
// 	}
// 	var buff bytes.Buffer
// 	if err := gameRowTemplate.Execute(&buff, data); err != nil {
// 		return "", err
// 	}
// 	return buff.String(), nil
// }

// func getGames(db *sqlx.DB) *gamesInfo {
// 	var tableRows []string
// 	games := []gameDBData{}
// 	err := db.Select(&games, "select id, DATE(game_date) as game_date, team_1_d, team_1_d_id, team_1_o, team_1_o_id, team_2_d, team_2_d_id, team_2_o, team_2_o_id, team_1_half, team_2_half, team_1_final, team_2_final from last_games where game_date > cast(current_timestamp() as date) + interval -14 day;")
// 	if err != nil {
// 		log.Fatal(err)
// 	} else if len(games[0].Date) == 0 {
// 		log.Fatal("Select returned nothing")
// 	} else {
// 		for _, game := range games {
// 			st, err := buildRow(game)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			tableRows = append(tableRows, st)
// 		}
// 	}
// 	return &gamesInfo{Games: template.HTML(strings.Join(tableRows, "\n"))}
// }

type gamesInfo struct {
	Games template.HTML
}

// RenderGamesPage : gets data from db to show last 10 games played
// func RenderGamesPage(db *sqlx.DB, r *http.Request) (template.HTML, error) {
// 	t, err := template.ParseFiles("webpage/games_view/games_template.html")
// 	if err != nil {
// 		return template.HTML(""), err
// 	}
// 	g := getGames(db)
// 	var buff bytes.Buffer
// 	if err = t.Execute(&buff, g); err != nil {
// 		return template.HTML(""), err
// 	}
// 	return template.HTML(buff.String()), nil
// }

type gameEditData struct {
	ID      int
	Date    string
	T1pd    int
	T1po    int
	T2pd    int
	T2po    int
	T1half  int
	T2half  int
	T1final int
	T2final int
}

// SaveGameEdit : saves a PUT request to alter a games data
func SaveGameEdit(env *util.Env, w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var t gameEditData
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&t)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			updateGame := `UPDATE games SET
							team_1_p1 = ?, team_1_p2 = ?, team_2_p1 = ?, team_2_p2 = ?,
							team_1_half = ?, team_2_half = ?, team_1_final = ?, team_2_final = ? 
							WHERE id = ?;`
			var res sql.Result
			if t.T1final > t.T2final {
				res, err = env.DB.Exec(updateGame, t.T1pd, t.T1po, t.T2pd, t.T2po, t.T1half, t.T2half, t.T1final, t.T2final, t.ID)
			} else {
				res, err = env.DB.Exec(updateGame, t.T2pd, t.T2po, t.T1pd, t.T1po, t.T2half, t.T1half, t.T2final, t.T1final, t.ID)
			}
			if err != nil {
				log.Print(err)
				http.Error(w, "SQL failed", http.StatusInternalServerError)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil || rowCnt != 1 {
				http.Error(w, "Nothing updated", http.StatusInternalServerError)
			}
		}
	} else {
		http.Error(w, "NOT ALLOWED", http.StatusMethodNotAllowed)
	}
}
