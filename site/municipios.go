package site

import (
    "appengine"
	"html/template"
    "net/http"
	"fmt"
	"model"
	"sess"
)

func init() {
    http.HandleFunc("/msp", municipios)
}

func municipios(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	_, ok := sess.IsSess(w, r, c)
	if ok {
		if entidad, err := model.GetEntidad(c, r.FormValue("CveEnt")); err == nil {
			if municipios, _ := entidad.GetMunicipios(c); err == nil {
				//Despliega municipios
				tpl, _ := template.New("Mun").Parse(OptionTpl)
				fmt.Fprintf(w, `<select name="DirMun" class="last" id="MunSelector" onchange="locateAddress();">`)
				for _, m := range *municipios {
					//fmt.Fprintf(w, "mun: %s, %s", r.FormValue("DirEnt"), m)
					// Ojo: ver porqué los repite con datos en blanco
					// El if es para corregir la bronca temporalmente
					if(m.Municipio != "") {
						if (m.CveMun == r.FormValue("CveMun")) {
							m.Selected = "selected"
						}
						tpl.Execute(w, m)
					}
				}
				fmt.Fprintf(w, `</select>`)
			}
		}
	}
	return
}

const OptionTpl = `<option value="{{.CveMun}}" {{if .Selected}}selected="{{.Selected}}"{{end}}>{{.Municipio}}</option>
`
