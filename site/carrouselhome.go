package site

import (
    "appengine"
    "appengine/datastore"
	"html/template"
	"math/rand"
    "net/http"
	"sortutil"
    "model"
	"time"
    "fmt"
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

const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('/spic?IdEmp={{.IdEmp}}')"></div></div>`
//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><img class="centerimg" src="/spic?IdEmp={{.IdEmp}}" /></div>`

func directorioLogos(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/*
		Loop para recorrer todas las empresas 
	*/
	prefix := r.FormValue("prefix")
	date := r.FormValue("bydate")
	name := r.FormValue("showname")
	all := r.FormValue("all")
	//o, _ := strconv.Atoi(r.FormValue("o"))
    //q := datastore.NewQuery("Cta").Order("FechaHora").Limit(n).Offset(o)
    q := datastore.NewQuery("Empresa")
	if all != "1" {
		q = q.Filter("Nombre >=", prefix).Filter("Nombre <", prefix+"\ufffd")
	}
	em, _ := q.Count(c)
	empresas := make([]model.Empresa, 0, em)
	if _, err := q.GetAll(c, &empresas); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return
		}
	}

	sortutil.AscByField(empresas, "Nombre")
	if date != "" {
		sortutil.AscByField(empresas, "FechaHora")
	}
	for k, _ := range empresas {
		imgq := datastore.NewQuery("EmpLogo").Filter("IdEmp =", empresas[k].IdEmp)
		for imgcur := imgq.Run(c); ; {
			var img model.Image
			_, err := imgcur.Next(&img)
			if err == datastore.Done  {
				break
			}
			var val string
			if name!="" {
				val = empresas[k].Nombre
			}
			if(img.Data != nil) {
				fmt.Fprintf(w, "<div class=\"imgCont\"><a href=\"/spic?IdEmp=%s\"><img src=\"/spic?IdEmp=%s\" />%s</a></div>\n", empresas[k].IdEmp, empresas[k].IdEmp, val)
			}
		}

	}
}

