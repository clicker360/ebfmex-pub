package site

import (
    "appengine"
	"html/template"
    "net/http"
	"sess"
)

func init() {
    http.HandleFunc("/prox", prox)
}

func prox(w http.ResponseWriter, r *http.Request) {
c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		proxTpl.Execute(w, s)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

var proxTpl = template.Must(template.ParseFiles("templates/proximamente.html"))
