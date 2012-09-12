package site

import (
    "net/http"
    "html/template"
)

func init() {
    http.HandleFunc("/sucursales", Sucursales)
}

func Sucursales(w http.ResponseWriter, r *http.Request){
	mapa := "Nada"
	sucursalesTpl.Execute(w, mapa)
}

var sucursalesTpl = template.Must(template.ParseFiles("templates/ListSucusales.html"))
