package home

import (
    "appengine"
    "appengine/datastore"
    "appengine/memcache"
	"html/template"
	"encoding/json"
	"math/rand"
    "net/http"
	"sortutil"
	"strings"
	"strconv"
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
    http.HandleFunc("/dirtexto", directorioTexto)
    http.HandleFunc("/wsdiremp", wsDirTexto)
    http.HandleFunc("/carr", carr)
    rand.Seed(time.Now().UnixNano())
}

/*
 * La idea es hacer 60 cachés al azar con un tiempo de vida de 30 min
 * Cada que se muere un memcache se genera otro carrousel al azar de logos
 */
func carr(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
    c := appengine.NewContext(r)
	var timetolive = 1800 //seconds
	var b []byte
	var nn int = 50 // tamaño del carrousel
	logos := make([]model.Image, 0, nn)
	hit := rand.Intn(60)
	if item, err := memcache.Get(c, "carr_"+strconv.Itoa(hit)); err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("ShortLogo")
		n, _ := q.Count(c)
		offset := 0;
		if(n > nn) {
			offset = rand.Intn(n-nn)
		} else {
			nn = n
		}
		q = q.Offset(offset).Limit(nn)
		if _, err := q.GetAll(c, &logos); err != nil {
			if err == datastore.ErrNoSuchEntity {
				return
			}
		}

		b, _ = json.Marshal(logos)
		item := &memcache.Item{
			Key:   "carr_"+strconv.Itoa(hit),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Memcache.Add carr_idoft : %v", err)
		}
		//c.Infof("memcache add carr_page : %v", strconv.Itoa(hit))
	} else {
		//c.Infof("memcache retrieve carr_page : %v", strconv.Itoa(hit))
		if err := json.Unmarshal(item.Value, &logos); err != nil {
			c.Errorf("Unmarshaling ShortLogo item: %v", err)
		}
		nn = len(logos)
	}

	tpl, _ := template.New("Carr").Parse(cajaTpl)
	tn := rand.Perm(nn)
	var ti response
	for i, _ := range tn {
		ti.IdEmp = logos[tn[i]].IdEmp
		ti.Name = logos[tn[i]].Name
		//b, _ := json.Marshal(ti)
		//w.Write(b)
		tpl.Execute(w, ti)
		if i >= nn  {
			break
		}
	}
}

func directorioTexto(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

	/*
		Loop para recorrer todas las empresas 
	*/
	now := time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
	prefixu := strings.ToUpper(r.FormValue("prefix"))
	ultimos := r.FormValue("ultimos")
    q := datastore.NewQuery("Empresa")
	if ultimos != "1" {
		q = q.Filter("Nombre >=", prefixu).Filter("Nombre <", prefixu+"\ufffd").Limit(400)
	} else {
		q = q.Filter("FechaHora >=", now.AddDate(0,0,-2)).Limit(400)
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


