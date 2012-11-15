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

type carrst struct {
	Name	string `json:"name"`
	Url		string `json:"url"`
}

type dirst struct {
	IdEmp	string `json:"id"`
	Name	string `json:"name"`
	Url		string `json:"url"`
	Num		int `json:"num"`
}

type qcount struct {
	lot	string `json:"count"`
}

type Paginador struct {
	Prefix string
	Pagina int
}

func init() {
    rand.Seed(time.Now().UnixNano())
    http.HandleFunc("/dirtexto", directorioTexto)
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
	var timetolive = 21600 //seconds
	var b []byte
	var nn int = 50 // tamaño del carrousel
	var logos [50]carrst
	hit := rand.Intn(50)
	cachename := "carr_"+strconv.Itoa(hit)
	if item, err := memcache.Get(c, cachename); err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("EmpLogo")
		n, _ := q.Count(c)
		offset := 0;
		if(n > nn) {
			offset = rand.Intn(n-nn)
		} else {
			nn = n
		}
		q = q.Offset(offset).Limit(nn)
		var ii int = 0
		for i := q.Run(c); ; {
			var e model.Image
			_, err := i.Next(&e)
			if err == datastore.Done {
				break
			}
			if e.IdEmp != "" {
				logos[ii].Name = e.Name
				logos[ii].Url = e.Sp4
				ii = ii+1
			}
		}

		/*
		if _, err := q.GetAll(c, &logos); err != nil {
			if err == datastore.ErrNoSuchEntity {
				return
			}
		}
		*/
		nn = len(logos)
		b, _ = json.Marshal(logos)
		item := &memcache.Item{
			Key:   cachename,
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("memcache.Add %v : %v", cachename, err)
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("Memcache.Add %v : %v", cachename, err)
			} else {
				c.Infof("memcached %v", cachename)
			}
		} else {
			c.Infof("memcached %v", cachename)
		}
	} else {
		if err := json.Unmarshal(item.Value, &logos); err != nil {
			c.Errorf("Memcache Unmarshalling %v : %v", cachename, err)
		}
		nn = len(logos)
	}

	tpl, _ := template.New("Carr").Parse(cajaTpl)
	tn := rand.Perm(nn)
	var ti carrst
	for i, _ := range tn {
		ti.Name = logos[tn[i]].Name
		ti.Url = strings.Replace(logos[tn[i]].Url, "s180", "s70",1)
		if ti.Url != "" {
			//b, _ := json.Marshal(ti)
			//w.Write(b)
			tpl.Execute(w, ti)
		}
		if i >= nn  {
			break
		}
	}
}

func directorioTexto(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	prefixu := strings.ToLower(r.FormValue("prefix"))
	ultimos := r.FormValue("ultimos")
	if ultimos == "" {
		if prefixu == "" || !model.ValidAlfa.MatchString(prefixu) || len(prefixu) > 1 {
			return
		}
	}

	/*
		Loop para recorrer todas las empresas 
	*/
	now := time.Now().Add(time.Duration(model.GMTADJ)*time.Second)
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	page -= 1
	const batch = 200
    q := datastore.NewQuery("EmpresaNm")
	var timetolive = 21600 //seconds
	if ultimos != "1" && prefixu !="" {
		var empresas []model.EmpresaNm
		var lot int
		cachename := "dirprefix_count_"+prefixu
		q = q.Filter("Nombre >=", prefixu).Filter("Nombre <", prefixu+"\ufffd").Order("Nombre")
		if item, err := memcache.Get(c, cachename); err == memcache.ErrCacheMiss {
			/*
			 * Se pagina ordenado alfabéticamente el resutlado de la búsqueda 
			 * y se guarda en Memcache
			 */
			lot, _ := q.Count(c)
			slot := strconv.Itoa(lot)
			item := &memcache.Item {
				Key:   cachename,
				Value: []byte(slot),
				Expiration: time.Duration(timetolive)*time.Second,
			}
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("memcache.Add %v : %v", cachename, err)
				if err := memcache.Set(c, item); err == memcache.ErrNotStored {
					c.Errorf("Memcache.Set %v : %v", cachename, err)
				} else {
					c.Infof("memcached %v", cachename)
				}
			} else {
				c.Infof("memcached %v", cachename)
			}
		} else {
			lot,_ = strconv.Atoi(string(item.Value))
		}

		pages := lot/batch
		if lot%batch > 0 {
			pages += 1
		}
		if pages > 1 {
			Paginas := make([]Paginador, pages)
			//c.Infof("lote: %d, paginas : %d", lot, pages)
			for i := 0; i < pages; i++ {
				Paginas[i].Prefix = prefixu
				Paginas[i].Pagina = i+1
				//c.Infof("pagina : %d", i)
			}
			tplp, _ := template.New("paginador").Parse(paginadorTpl)
			tplp.Execute(w, Paginas)
		}
		cachename = "dirprefix_"+prefixu+"_"+strconv.Itoa(page)
		if item, err := memcache.Get(c, cachename); err == memcache.ErrCacheMiss {
			offset := batch * page
			q = q.Offset(offset).Limit(batch)
			if _, err := q.GetAll(c, &empresas); err != nil {
				return
			}

			b, _ := json.Marshal(empresas)
			item := &memcache.Item {
				Key:   cachename,
				Value: b,
				Expiration: time.Duration(timetolive)*time.Second,
			}
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("memcache.Add %v : %v", cachename, err)
				if err := memcache.Set(c, item); err == memcache.ErrNotStored {
					c.Errorf("Memcache.Set %v : %v", cachename, err)
				} else {
					c.Infof("memcached %v", cachename)
				}
			} else {
				c.Infof("memcached %v", cachename)
			}
		} else {
			if err := json.Unmarshal(item.Value, &empresas); err != nil {
				c.Errorf("Memcache Unmarshalling %v : %v", cachename, err)
			}
		}

		sortutil.CiAscByField(empresas, "Nombre")
		var ti dirst
		var tictac int
		//var repetido string
		tictac = 1
		for k, _ := range empresas {
			tpl, _ := template.New("pagina").Parse(empresaTpl)
			ti.Num = tictac
			ti.IdEmp = empresas[k].IdEmp
			ti.Name = strings.Title(empresas[k].Nombre)
			//if repetido != ti.Name {
				if tictac != 1 {
					tictac = 1
				} else {
					tictac = 2
				}
				//repetido = ti.Name
				tpl.Execute(w, ti)
			//}
		}
	} else {
		prefixu = ""
		var empresas []model.EmpresaNm
		cachename := "dirprefix_ultimos"
		if item, err := memcache.Get(c, cachename); err == memcache.ErrCacheMiss {
			q = datastore.NewQuery("Empresa").Filter("FechaHora >=", now.AddDate(0,0,-2)).Limit(300)
			var empresas []model.Empresa
			if _, err := q.GetAll(c, &empresas); err != nil {
				return
			}
			b, _ := json.Marshal(empresas)
			item := &memcache.Item {
				Key:   cachename,
				Value: b,
				Expiration: time.Duration(timetolive)*time.Second,
			}
			if err := memcache.Add(c, item); err == memcache.ErrNotStored {
				c.Errorf("memcache.Add %v : %v", cachename, err)
				if err := memcache.Set(c, item); err == memcache.ErrNotStored {
					c.Errorf("memcache.Set %v : %v", cachename, err)
				} else {
					c.Infof("memcached %v", cachename)
				}
			} else {
				c.Infof("memcached %v", cachename)
			}
		} else {
			if err := json.Unmarshal(item.Value, &empresas); err != nil {
				c.Errorf("Memcache Unmarshalling %v : %v", cachename, err)
			}
		}

		sortutil.CiAscByField(empresas, "Nombre")
		var ti dirst
		var tictac int
		//var repetido string
		tictac = 1
		for k, _ := range empresas {
			tpl, _ := template.New("pagina").Parse(empresaTpl)
			ti.Num = tictac
			ti.IdEmp = empresas[k].IdEmp
			ti.Name = strings.Title(strings.ToLower(empresas[k].Nombre))
			//if repetido != ti.Name {
				if tictac != 1 {
					tictac = 1
				} else {
					tictac = 2
				}
				//repetido = ti.Name
				tpl.Execute(w, ti)
			//}
		}
	}
}

//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('/spic?IdEmp={{.IdEmp}}')"></div></div>`
const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><div class="centerimg" style="background-image:url('{{.Url}}')"></div></div>`
const empresaTpl = `<div class="gridsubRow bg-Gry{{.Num}}"><a href="http://www.elbuenfin.org/micrositio.html?id={{.IdEmp}}" target="_blank">{{.Name}}</a></div>`
const paginadorTpl = `<div class="pagination-H"><ul id="letters">{{range .}}<li><a href="#" class="letter" prfx="{{.Prefix}}" onclick="javascript:paginar({{.Pagina}});"> {{.Pagina}} </a></li>{{end}}</ul></div>`
//const paginadorTpl = `<div>{{range .}}<a href="javascript:pager({{.Prefix}}, {{.Pagina}});"> {{.Pagina}} </a>{{end}}</div>`
//const cajaTpl = `<div class="cajaBlanca" title="{{.Name}}"><img class="centerimg" src="/spic?IdEmp={{.IdEmp}}" /></div>`
