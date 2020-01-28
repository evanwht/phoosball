package gopages

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"log"
	"io/ioutil"

	"github.com/evanwht1/phoosball/util"
)

// AddNewPlayer : adds a new player to the database
func AddNewPlayer(db *sql.DB, w http.ResponseWriter, r *http.Request) *util.Page {
	r.ParseForm()
	if len(r.PostForm) > 0 {
		p, err := template.ParseFiles("webpage/account_template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var buff bytes.Buffer
		if err = p.Execute(&buff, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return &util.Page{Title: "Account", Body: template.HTML(buff.String())}
	}
	return nil
}

// RenderPlayerPage : renders the game input form page with correct data
func RenderPlayerPage(db *sql.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	r.ParseForm()
	var AlertMessage template.HTML
	if len(r.PostForm) > 0 {
		fail := false
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
			fail = true
		} else {
			stmt, err := db.Prepare(`INSERT INTO players (name, display_name, email) VALUES (?, ?, ?);`)
			if err != nil {
				log.Println(err)
				fail = true
			}
			res, err := stmt.Exec(r.PostFormValue("t1_p1"), r.PostFormValue("t1_p2"), r.PostFormValue("t2_p1"), r.PostFormValue("t2_p2"))
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
