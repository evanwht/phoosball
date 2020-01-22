package gopages

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func buildRow(date string, t1Pd string, t1Po string, t2Pd string, t2Po string, t1Half string, t2Half string, t1End string, t2End string) string {
	var (
		team1Class string
		team2Class string
	)
	if t1End > t2End {
		team1Class = "success"
		team2Class = "danger"
	} else {
		team1Class = "danger"
		team2Class = "success"
	}
	// TODO change to template
	return `<tr>
				<th scope="row">` + date + `</th>
				<td class="text-` + team1Class + `">` + t1Pd + " - " + t1Po + `</td>
				<td class="text-` + team2Class + `">` + t2Pd + " - " + t2Po + `</td>
				<td>` + t1Half + " - " + t2Half + `</td>
				<td>` + t1End + " - " + t2End + `</td>
				<td><button type="button"class="btn btn-outline-warning">Edit</button></td>
			</tr>`
}

func getGames(db *sql.DB) *gamesInfo {
	var (
		date       string
		t1Pd      string
		t1Po      string
		t2Pd      string
		t2Po      string
		t1Half    string
		t2Half    string
		t1End     string
		t2End     string
		tableRows []string
	)
	rows, err := db.Query("select * from last_games;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&date, &t1Pd, &t1Po, &t2Pd, &t2Po, &t1Half, &t2Half, &t1End, &t2End)
			if err != nil {
				log.Fatal(err)
			} else {
				tableRows = append(tableRows, buildRow(date[:11], t1Pd, t1Po, t2Pd, t2Po, t1Half, t2Half, t1End, t2End))
			}
		}
		rows.Close()
	}
	return &gamesInfo{Games: template.HTML(strings.Join(tableRows, "\n"))}
}

type gamesInfo struct {
	Games template.HTML
}

// RenderPage : gets data from db to show last 10 games played
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
