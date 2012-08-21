package site

import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "sess"
    "html/template"
    "time"
)

type Cta struct {
        Folio                   int32
        Nombre                  string
        Apellidos               string
        Puesto                  string
        Email                   string
        EmailAlt                string
        Pass                    string
        Tel                             string
        Cel                             string
        FechaHora               time.Time
        UsuarioInt              string
        CodigoCfm               string
        Status                  bool
}

func init() {
    http.HandleFunc("/registro-export", registroExport)
    http.HandleFunc("/registros.csv", registroCsv)
}

func registroExport(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        if _, ok := sess.IsSess(w, r, c); !ok {
                tc := make(map[string]interface{})
                exportTpl.ExecuteTemplate(w, "cta", tc)
                return
        } else {
		return
        }
}

func registroCsv(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Cta").Order("FechaHora")
    regdata := make([]Cta,0,10)
    if _, err := q.GetAll(c, &regdata); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
        w.Header().Set("Content-Type", "text/csv")
    if err := registrosCsvTpl.Execute(w, regdata); err != nil {
	return
    }
}

var exportTpl = template.Must(template.ParseFiles("templates/export.html"))
var registrosCsvTpl = template.Must(template.ParseFiles("templates/registros.csv"))
