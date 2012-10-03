package site

import (
    "appengine"
    "appengine/datastore"
	"html/template"
	"math/rand"
    "net/http"
	"sortutil"
	"strings"
    "model"
	"time"
)

type cimage struct {
	IdEmp	string `json:"id"`
	Name	string `json:"name"`
}

func init() {
    rand.Seed(time.Now().UnixNano())
    http.HandleFunc("/dirlogos", directorioLogos)
    http.HandleFunc("/carr", carr)
}

func carr(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
    c := appengine.NewContext(r)
    q := datastore.NewQuery("ShortLogo")
	n, _ := q.Count(c)
	logos := make([]model.Image, 0, n)
	if _, err := q.GetAll(c, &logos); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return
		}
	}
	tpl, _ := template.New("Carr").Parse(cajaTpl)
	tn := rand.Perm(n)
	var ti cimage
	for i, _ := range tn {
		ti.IdEmp = logos[tn[i]].IdEmp
		ti.Name = logos[tn[i]].Name
		//b, _ := json.Marshal(ti)
		//w.Write(b)
		tpl.Execute(w, ti)
		if i > 49  {
			break
		}
	}
}

func directorioLogos(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/*
		Loop para recorrer todas las empresas 
	*/
	prefixu := strings.ToUpper(r.FormValue("prefix"))
	//prefixl := strings.ToLower(r.FormValue("prefix"))
	//date := r.FormValue("bydate")
	ultimos := r.FormValue("ultimos")
    q := datastore.NewQuery("Empresa")
	if ultimos != "1" {
		q = q.Filter("Nombre >=", prefixu).Filter("Nombre <", prefixu+"\ufffd")
	} else {
		q = q.Filter("FechaHora >=", time.Now().AddDate(0,0,-2))
	}
	em, _ := q.Count(c)
	empresas := make([]model.Empresa, 0, em)
	if _, err := q.GetAll(c, &empresas); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return
		}
	}

	sortutil.AscByField(empresas, "Nombre")
	/*
	if date != "" {
		sortutil.AscByField(empresas, "FechaHora")
	}
	*/
	for k, _ := range empresas {
		imgq := datastore.NewQuery("EmpLogo").Filter("IdEmp =", empresas[k].IdEmp)
		for imgcur := imgq.Run(c); ; {
			var img model.Image
			_, err := imgcur.Next(&img)
			if err == datastore.Done  {
				break
			}
			tpl, _ := template.New("Carr").Parse(cajaTpl)
			var ti cimage
			if(img.Data != nil) {
				ti.IdEmp = empresas[k].IdEmp
				ti.Name = empresas[k].Nombre
				tpl.Execute(w, ti)
			}
		}

	}
}

const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('/spic?IdEmp={{.IdEmp}}')"></div></div>`
//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><img class="centerimg" src="/spic?IdEmp={{.IdEmp}}" /></div>`

