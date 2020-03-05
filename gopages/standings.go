package gopages

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func standingsRow(place int, cur standing) string {
	return `<tr>
				<th scope="row">` + strconv.Itoa(place) + `</th>
				<td>` + cur.Name + `</td>
				<td>` + strconv.Itoa(cur.Wins) + `</td>
				<td>` + strconv.Itoa(cur.Losses) + `</td>
				<td>` + fmt.Sprintf("%01.2f", cur.Perc) + `</td>
			</tr>`
}

type standing struct {
	Name   string
	Wins   int
	Losses int
	Perc   float32
}

func getStandings(db *sqlx.DB) *standingsInfo {
	var tableRows []string
	standings := []standing{}
	err := db.Select(&standings, "select *, (wins / (wins+losses)) perc from overall_standings order by perc desc;")
	if err != nil {
		log.Fatal(err)
	} else {
		for i, standing := range standings {
			tableRows = append(tableRows, standingsRow(i, standing))
		}
	}
	return &standingsInfo{Standings: template.HTML(strings.Join(tableRows, "\n"))}
}

type standingsInfo struct {
	Standings template.HTML
}

// RenderStandingsPage : gets data from db to show current standings of known players
func RenderStandingsPage(db *sqlx.DB, r *http.Request) (template.HTML, error) {
	t, err := template.ParseFiles("webpage/index_template.html")
	if err != nil {
		return template.HTML(""), err
	}
	g := getStandings(db)
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(buff.String()), nil
}
