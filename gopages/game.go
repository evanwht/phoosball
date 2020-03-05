package gopages

import (
	"bytes"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/evanwht1/phoosball/util"
)

// CreatePlayerOptions : creates select options for each player in the db
func CreatePlayerOptions(db *sqlx.DB) template.HTML {
	rows := GetAllPlayers(db)
	return template.HTML(strings.Join(playersToOptions(rows), "\n"))
}

func playersToOptions(players []Player) []string {
	var options []string
	for _, player := range players {
		options = append(options, util.HTMLOption(strconv.Itoa(player.ID), player.Name+" ("+player.NickName+")"))
	}
	return options
}

type goalType struct {
	ID   string
	Name string
}

// CreateGoalOptions : creates select options for each goal type in the db
func CreateGoalOptions(db *sqlx.DB) template.HTML {
	var events []string
	goals := []goalType{}
	err := db.Select(&goals, "select * from event_types;")
	if err != nil {
		log.Fatal(err)
	} else {
		for _, goal := range goals {
			events = append(events, util.HTMLOption(goal.ID, goal.Name))
		}
	}
	return template.HTML(strings.Join(events, "\n"))
}

type gameInfo struct {
	PlayerOptions template.HTML
	GoalOptions   template.HTML
	Alert         template.HTML
}

var fallBackAlert = "<div><p class=\"text-danger\">Unknown failure. Contanct admin</p></div>"

// RenderGamePage : renders the game input form page with correct data
func RenderGamePage(db *sqlx.DB, r *http.Request) (template.HTML, error) {
	r.ParseForm()
	var AlertMessage template.HTML
	if len(r.PostForm) > 0 {
		// User has submitted a game page data. try to insert in to db or return error message
		fail := false
		insertSQL := `INSERT INTO games 
			(team_1_p1, team_1_p2, team_2_p1, team_2_p2,
			team_1_half, team_2_half, team_1_final, team_2_final) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
		t1Final, err1 := strconv.Atoi(r.PostFormValue("t1_final"))
		t2Final, err2 := strconv.Atoi(r.PostFormValue("t2_final"))
		if err1 != nil || err2 != nil {
			log.Println("could not parse final scores from new game")
			fail = true
		} else {
			var (
				res sql.Result
				err error
			)
			if t1Final > t2Final {
				res, err = db.Exec(insertSQL, r.PostFormValue("t1_p1"), r.PostFormValue("t1_p2"),
					r.PostFormValue("t2_p1"), r.PostFormValue("t2_p2"),
					r.PostFormValue("t1_half"), r.PostFormValue("t2_half"),
					r.PostFormValue("t1_final"), r.PostFormValue("t2_final"))
			} else {
				res, err = db.Exec(insertSQL, r.PostFormValue("t2_p1"), r.PostFormValue("t2_p2"),
					r.PostFormValue("t1_p1"), r.PostFormValue("t1_p2"),
					r.PostFormValue("t2_half"), r.PostFormValue("t1_half"),
					r.PostFormValue("t2_final"), r.PostFormValue("t1_final"))
			}
			if err != nil {
				log.Println(err)
				fail = true
			}
			rowCnt, err := res.RowsAffected()
			if err != nil || rowCnt <= 0 {
				fail = true
			}
		}
		// show message alert
		if fail {
			b, err := ioutil.ReadFile("webpage/game_input/fail_alert.html")
			if err != nil {
				AlertMessage = template.HTML(fallBackAlert)
			}
			AlertMessage = template.HTML(string(b))
		} else {
			b, err := ioutil.ReadFile("webpage/game_input/success_alert.html")
			if err != nil {
				AlertMessage = template.HTML(fallBackAlert)
			}
			AlertMessage = template.HTML(string(b))
		}
	}

	opts := CreatePlayerOptions(db)
	g := &gameInfo{PlayerOptions: opts, Alert: AlertMessage}

	t, err := template.ParseFiles("webpage/game_input/game_template.html")
	if err != nil {
		return template.HTML(""), err
	}
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(buff.String()), nil
}
