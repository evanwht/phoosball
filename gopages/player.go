package gopages

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"

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
