package gopages

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func standingsRow(place int, name string, wins int, losses int, perc float32) string {
	return `<tr>
				<th scope="row">` + strconv.Itoa(place) + `</th>
				<td>` + name + `</td>
				<td>` + strconv.Itoa(wins) + `</td>
				<td>` + strconv.Itoa(losses) + `</td>
				<td>` + fmt.Sprintf("%01.2f", perc) + `</td>
			</tr>`
}

func getStandings(db *sql.DB) *standingsInfo {
	var (
		name      string
		wins      int
		losses    int
		perc      float32
		tableRows []string
	)
	rows, err := db.Query("select *, (wins / (wins+losses)) perc from overall_standings order by perc desc;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		i := 1
		for rows.Next() {
			err := rows.Scan(&name, &wins, &losses, &perc)
			if err != nil {
				log.Fatal(err)
			} else {
				tableRows = append(tableRows, standingsRow(i, name, wins, losses, perc))
				i++
			}
		}
		rows.Close()
	}
	return &standingsInfo{Standings: template.HTML(strings.Join(tableRows, "\n"))}
}

type standingsInfo struct {
	Standings template.HTML
}

// RenderStandingsPage : gets data from db to show current standings of known players
func RenderStandingsPage(db *sql.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	t, err := template.ParseFiles("webpage/index_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	g := getStandings(db)
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	return template.HTML(buff.String()), nil
}
