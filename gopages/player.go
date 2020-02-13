package gopages

import (
	"bytes"
	"github.com/jmoiron/sqlx"
	"html/template"
	"net/http"
	"log"
	"io/ioutil"

	"github.com/evanwht1/phoosball/util"
)

// GetAllPlayers : gets all selectable players from the db
func GetAllPlayers(db *sqlx.DB) []Player {
	var players []Player
	rows, err := db.Query("select id, name, display_name from players;")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var p Player
			err := rows.Scan(&p.ID, &p.Name, &p.NickName)
			if err != nil {
				log.Fatal(err)
			} else {
				players = append(players, p)
			}
		}
		rows.Close()
	}
	return players
}

// AddNewPlayer : adds a new player to the database
func getAccountPage(db *sqlx.DB, id int, w http.ResponseWriter, r *http.Request) *util.Page {
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

type playerInfo struct {
	Alert template.HTML
}

// RenderPlayerPage : renders the game input form page with correct data
func RenderPlayerPage(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	r.ParseForm()
	var AlertMessage template.HTML
	if len(r.PostForm) > 0 {
		fail := false
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
			fail = true
		} else {
			stmt, err := tx.Prepare(`INSERT INTO players (name, display_name, email) VALUES (?, ?, ?);`)
			if err != nil {
				log.Println(err)
				fail = true
			} else {
				res, err := stmt.Exec(r.PostFormValue("firstName") + " " + r.PostFormValue("lastName"), r.PostFormValue("nickName"), r.PostFormValue("email"))
				if err != nil {
					log.Println(err)
					fail = true
				} else {
					lastID, err := res.LastInsertId()
					if err != nil || lastID <= 0 {
						fail = true
					}
					rowCnt, err := res.RowsAffected()
					if err != nil || rowCnt <= 0 {
						fail = true
					}
				}
			}
			if fail {
				tx.Rollback()
			} else {
				tx.Commit()
			}
			stmt.Close()
		}
		// show message alert
		if fail {
			b, err := ioutil.ReadFile("webpage/player_input/fail_alert.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				AlertMessage =  template.HTML(fallBackAlert)
			}
			AlertMessage =  template.HTML(string(b))
		}
		b, err := ioutil.ReadFile("webpage/player_input/success_alert.html")
		if err != nil {
			AlertMessage =  template.HTML(fallBackAlert)
		}
		AlertMessage = template.HTML(string(b))
	}

	g := &playerInfo{Alert: AlertMessage}

	t, err := template.ParseFiles("webpage/player_input/player_template.html")
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
