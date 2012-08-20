package site

import (
    "appengine"
    "net/http"
)

func init() {
    http.HandleFunc("/prox", prox)
}

func prox(w http.ResponseWriter, r *http.Request) {
c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		proxTpl.Execute(w, s)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
