package site

import (
    "net/http"
    "html/template"
)

func init() {
    http.HandleFunc("/mapmkta", MapMkta)
}

func MapMkta(w http.ResponseWriter, r *http.Request){
	mapa := "Nada"
	mapMktaTpl.Execute(w, mapa)
}

var mapMktaTpl = template.Must(template.ParseFiles("templates/mapMkta.html"))
