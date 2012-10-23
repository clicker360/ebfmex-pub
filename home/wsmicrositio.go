package home

import (
    "appengine"
    "appengine/memcache"
	"encoding/json"
    "net/http"
    "model"
    "time"
)

type micrositio struct {
	IdEmp		string	`json:"idemp"`
	IdImg		string	`json:"idimg"`
	Name		string	`json:"name"`
	Desc		string	`json:"desc"`
	Url         string	`json:"url"`
	Sp1         string	`json:"facebook"`
	Sp2         string	`json:"twitter"`
}

func init() {
    http.HandleFunc("/wsmicrositio", ShowMicrositio)
}

func ShowMicrositio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var timetolive = 1800 //seconds
	c := appengine.NewContext(r)
	var b []byte
	var m micrositio
	if item, err := memcache.Get(c, "m_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
		if imgo := model.GetLogo(c, r.FormValue("id")); imgo != nil {
			m.IdEmp = imgo.IdEmp
			m.IdImg = imgo.IdImg
			m.Name = imgo.Name
			m.Desc = imgo.Desc
			m.Url = imgo.Url
			m.Sp1 = imgo.Sp1
			m.Sp2 = imgo.Sp2
		}

		b, _ = json.Marshal(m)
		item := &memcache.Item{
			Key:   "m_"+r.FormValue("id"),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Error memcache.Add m_idoft : %v", err)
			w.Write(b)
		}
	} else {
		//c.Infof("memcache retrieve m_idoft : %v", r.FormValue("id"))
		b = item.Value
	}
	w.Write(b)
}

