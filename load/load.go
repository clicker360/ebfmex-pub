package load

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "net/http"
	"fmt"
)

type Organismo struct {
	Siglas		string
	Nombre		string
	Selected	string
}

func init() {
	http.HandleFunc("/loado", loadOrg)
}

func loadOrg(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
    if u := user.Current(c); u != nil {
		o := []Organismo {
			{"ABM", "Asociación de Bancos de México", ""},
			{"AMIB","Asociación Mexicana de Intermediarios Bursátiles", ""},
			{"AMIS", "Asociación Mexicana de Instituciones de Seguros", ""},
			{"AMPICI", "Asociación Mexicana de Internet", ""},
			{"ANTAD", "Asociación Nacional de Tiendas de Autoservicio y Departamentales", ""},
			{"CANACINTRA", "Cámara Nacional de la Industria de Transformación", ""},
			{"CNA", "Consejo Mexicano de Hombres de Negocios", ""},
			{"COMCE", "Consejo Empresarial Mexicano de Comercio Exterior, Inversión y Tecnología", ""},
			{"CONCANACO/CONCAMIN", "Confederación de Cámaras Nacionales de Comercio, Servicio y Turismo", ""},
			{"COPARMEX", "Confederación Patronal de la República Mexicana", ""},
			{"OTRO", "Otro Organismo", ""},
		}
		for _, e := range o {
			fmt.Fprintf(w, "Organismo: %d, %d", e.Siglas, e.Nombre)
			_, err := datastore.Put(c, datastore.NewKey(c, "Organismo", e.Nombre, 0, nil), &e)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
