package gopages

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
)

func createGameInfo(db *sql.DB) *gameInfo {
	var (
		id          int
		displayName string
		name        string
		names       []string
		events      []string
	)
	// get user info
	rows, err := db.Query("select id, name, display_name from players;")
	if err != nil {
		// do nothing
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &displayName)
			if err != nil {
				displayName = ""
			}
			names = append(names, "\""+displayName+" ("+name+")\",")
		}
		names = append(names, "\"New Player\"")
	}

	// get event type info
	rows, err = db.Query("select * from event_types;")
	if err != nil {
		// do nothing
	} else {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err == nil {
				events = append(events, "\""+name+"\",")
			} else {
				// load page with error message
			}
		}
		events = append(events, "\"New Type\"")
	}
	return &gameInfo{Players: names, GoalTypes: events}
}

type gameInfo struct {
	Players      []string
	GoalTypes    []string
	AlertMessage template.HTML
}

// RenderGamePage : renders the game input form page with correct data
func RenderGamePage(db *sql.DB, w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	r.ParseForm()
	var AlertMessage template.HTML
	if len(r.PostForm) > 0 {
		// User has submitted a previous game page data. try to insert in to db or return error message
		// TODO insert into db
		if (1 == 1) {
			AlertMessage = template.HTML("<div class=\"alert alert-success\" role=\"alert\"><h4 class=\"alert-heading\">Well done!</h4><p>Aww yeah, you successfully read this important alert message.</p><hr><p class=\"mb-0\">Whenever you need to, be sure to use margin utilities to keep things nice and tidy.</p></div>")
		}
	}

	g := createGameInfo(db)
	g.AlertMessage = AlertMessage

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
