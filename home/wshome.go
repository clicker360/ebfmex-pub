package home

import (
    "appengine"
    "appengine/datastore"
	"html/template"
	"encoding/json"
	"math/rand"
    "net/http"
	"sortutil"
	"strings"
    "model"
	"time"
)

type response struct {
	IdEmp	string `json:"id"`
	Name	string `json:"name"`
	Num		int `json:"num"`
}

func init() {
    rand.Seed(time.Now().UnixNano())
    http.HandleFunc("/dirlogos", directorioLogos)
    http.HandleFunc("/dirtexto", directorioTexto)
    http.HandleFunc("/wsdiremp", wsDirTexto)
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
	var ti response
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
	var ti response
	for k, _ := range empresas {
		imgq := datastore.NewQuery("EmpLogo").Filter("IdEmp =", empresas[k].IdEmp)
		for imgcur := imgq.Run(c); ; {
			var img model.Image
			_, err := imgcur.Next(&img)
			if err == datastore.Done  {
				break
			}
			tpl, _ := template.New("Carr").Parse(cajaTpl)
			if(img.Data != nil) {
				ti.IdEmp = empresas[k].IdEmp
				ti.Name = empresas[k].Nombre
				tpl.Execute(w, ti)
			}
		}

	}
}

func directorioTexto(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/*
		Loop para recorrer todas las empresas 
	*/
	prefixu := strings.ToUpper(r.FormValue("prefix"))
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
	var ti response
	var tictac int
	tictac = 1
	for k, _ := range empresas {
		if tictac != 1 {
			tictac = 1
		} else {
			tictac = 2
		}
		tpl, _ := template.New("Carr").Parse(empresaTpl)
		ti.Num = tictac
		ti.IdEmp = empresas[k].IdEmp
		ti.Name = empresas[k].Nombre
		tpl.Execute(w, ti)
	}
}

const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('/spic?IdEmp={{.IdEmp}}')"></div></div>`
const empresaTpl = `<div class="gridsubRow bg-Gry{{.Num}}">{{.Name}}</div>`
//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><img class="centerimg" src="/spic?IdEmp={{.IdEmp}}" /></div>`

type WsEmpresa struct{
	Id		string `json:"id"`
	Empresa	string `json:"empresa"`
	Url		string `json:"url"`
}

func wsDirTexto(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	prefixu := strings.ToUpper(r.FormValue("prefix"))
    q := datastore.NewQuery("Empresa").Filter("Nombre >=", prefixu).Filter("Nombre <", prefixu+"\ufffd") //.Filter("Status =", true)
	em, _ := q.Count(c)
	empresas := make([]model.Empresa, em, em)
	if _, err := q.GetAll(c, empresas); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return
		}
	}

	var b []byte
	wsout := make([]WsEmpresa, em, em)
	sortutil.AscByField(empresas, "Nombre")
	for i, _ := range empresas {
		wsout[i].Id = empresas[i].IdEmp
		wsout[i].Empresa = empresas[i].Nombre
		wsout[i].Url = empresas[i].Url
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ = json.Marshal(wsout)
	w.Write(b)
}


