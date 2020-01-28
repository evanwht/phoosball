package gopages

import (
	"bytes"
	"database/sql"
	"html/template"
	"io/ioutil"
	"strconv"
	"log"
	"net/http"
	"strings"

	"github.com/evanwht1/phoosball/util"
)

func createPlayerOptions(db *sql.DB) template.HTML {
	var (
		id          string
		displayName string
		name        string
		options     []string
	)
	rows, err := db.Query("select id, name, display_name from players;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &displayName)
			if err != nil {
				log.Fatal(err)
			} else {
				options = append(options, util.HTMLOption(id, displayName+" ("+name+")"))
			}
		}
		rows.Close()
	}
	return template.HTML(strings.Join(options, "\n"))
}

func createGoalTypes(db *sql.DB) template.HTML {
	var (
		id     string
		name   string
		events []string
	)
	rows, err := db.Query("select * from event_types;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			} else {
				events = append(events, util.HTMLOption(id, name))
			}
		}
		rows.Close()
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
func RenderGamePage(db *sql.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	r.ParseForm()
	var AlertMessage template.HTML
	if len(r.PostForm) > 0 {
		// User has submitted a game page data. try to insert in to db or return error message
		fail := false
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
			fail = true
		} else {
			stmt, err := db.Prepare(
				`INSERT INTO games 
				(team_1_p1, team_1_p2, team_2_p1, team_2_p2,
				team_1_half, team_2_half, team_1_final, team_2_final) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?);`)
			if err != nil {
				log.Println(err)
				fail = true
			}
			t1Final, err := strconv.Atoi(r.PostFormValue("t1_final"))
			if err != nil {
				log.Println(err)
				fail = true
			}
			t2Final, err := strconv.Atoi(r.PostFormValue("t2_final"))
			if err != nil {
				log.Println(err)
				fail = true
			}
			var res sql.Result
			if (t1Final > t2Final) {
				res, err = stmt.Exec(r.PostFormValue("t1_p1"), r.PostFormValue("t1_p2"),
					r.PostFormValue("t2_p1"), r.PostFormValue("t2_p2"),
					r.PostFormValue("t1_half"), r.PostFormValue("t2_half"),
					r.PostFormValue("t1_final"), r.PostFormValue("t2_final"))
			} else {
				res, err = stmt.Exec(r.PostFormValue("t2_p1"), r.PostFormValue("t2_p2"),
					r.PostFormValue("t1_p1"), r.PostFormValue("t1_p2"),
					r.PostFormValue("t2_half"), r.PostFormValue("t1_half"),
					r.PostFormValue("t2_final"), r.PostFormValue("t1_final"))
			}
			if err != nil {
				log.Println(err)
				fail = true
			}
			lastID, err := res.LastInsertId()
			if err != nil || lastID <= 0 {
				fail = true
			}
			rowCnt, err := res.RowsAffected()
			if err != nil || rowCnt <= 0 {
				fail = true
			}
			if fail {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}
		// show message alert
		if fail {
			b, err := ioutil.ReadFile("webpage/game_input/fail_alert.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				AlertMessage =  template.HTML(fallBackAlert)
			}
			AlertMessage =  template.HTML(string(b))
		}
		b, err := ioutil.ReadFile("webpage/game_input/success_alert.html")
		if err != nil {
			AlertMessage =  template.HTML(fallBackAlert)
		}
		AlertMessage = template.HTML(string(b))
	}

	opts := createPlayerOptions(db)
	g := &gameInfo{PlayerOptions: opts, Alert: AlertMessage}

	t, err := template.ParseFiles("webpage/game_template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	var buff bytes.Buffer
	if err = t.Execute(&buff, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return template.HTML(""), err
	}
	return template.HTML(buff.String()), nil
}
