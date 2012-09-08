package site

import (
    "net/http"
    "html/template"
)

func init() {
    http.HandleFunc("/mapa", Mapa)
}

func Mapa(w http.ResponseWriter, r *http.Request){
	mapa := "Nada"
	mapTpl.ExecuteTemplate(w, "Mapa", mapa)
}

var mapTpl = template.Must(template.ParseFiles("templates/map.html"))
