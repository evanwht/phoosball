package gopages

import (
	"bytes"
	"fmt"
	"database/sql"
	"html/template"
	"log"
	"strconv"
	"net/http"
	"strings"
)

func standingsRow(place int, name string, wins int, losses int) string {
	// TODO change to template
	// strconv.FormatFloat((float64(wins) / float64(wins+losses)), 'f', 0, 32)
	perc := fmt.Sprintf("%01.2f", (float64(wins) / float64(wins+losses)))
	return `<tr>
				<th scope="row">` + strconv.Itoa(place) + `</th>
				<td>` + name + `</td>
				<td>` + strconv.Itoa(wins) + `</td>
				<td>` + strconv.Itoa(losses) + `</td>
				<td>` + perc + `</td>
			</tr>`
}

func getStandings(db *sql.DB) *standingsInfo {
	var (
		name      string
		wins      int
		losses      int
		tableRows []string
	)
	rows, err := db.Query("select * from overall_standings;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		i := 1
		for rows.Next() {
			err := rows.Scan(&name, &wins, &losses)
			if err != nil {
				log.Fatal(err)
			} else {
				tableRows = append(tableRows, standingsRow(i, name, wins, losses))
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
