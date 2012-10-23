package home

import (
    "appengine"
    "appengine/memcache"
	"encoding/json"
    "net/http"
	"sortutil"
    "model"
    "time"
)

type wsoferta struct {
	IdEmp       string	`json:"idemp"`
	IdOft		string	`json:"idoft"`
	Oferta		string	`json:"oferta"`
	Descripcion	string	`json:"descripcion"`
	Enlinea		bool	`json:"enlinea"`
	Url			string	`json:"url"`
	BlobKey	appengine.BlobKey `json:"imgurl"`
}

func init() {
    http.HandleFunc("/wsofxe", ShowEmpOfertas)
}

func ShowEmpOfertas(w http.ResponseWriter, r *http.Request) {
	var timetolive = 1800 //seconds
	c := appengine.NewContext(r)
	var b []byte
	if item, err := memcache.Get(c, "ofxe_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
		emofs := model.ListOf(c, r.FormValue("id"))
		wsofs := make([]wsoferta, len(*emofs), cap(*emofs))
		for i,v:= range *emofs {
			wsofs[i].IdEmp = v.IdEmp
			wsofs[i].IdOft = v.IdOft
			wsofs[i].Oferta = v.Oferta
			wsofs[i].Descripcion = v.Descripcion
			wsofs[i].Enlinea = v.Enlinea
			wsofs[i].Url = v.Url
			wsofs[i].BlobKey = v.BlobKey
		}
		sortutil.AscByField(wsofs, "Oferta")
		b, _ = json.Marshal(wsofs)
		item := &memcache.Item{
			Key:   "ofxe_"+r.FormValue("id"),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Error memcache.Add ofxe_idemp : %v", err)
		}
	} else {
		//c.Infof("memcache retrieve sucs_idemp : %v", r.FormValue("id"))
		b = item.Value
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
