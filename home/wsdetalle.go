package home

import (
    "appengine"
    "appengine/memcache"
	"encoding/json"
    "net/http"
    "model"
    "time"
)

type detalle struct {
	IdEmp		string		`json:"idemp"`
	IdOft		string		`json:"idoft"`
	IdCat		int			`json:"idcat"`
	Oferta		string		`json:"oferta"`
	Empresa		string		`json:"empresa"`
	Descripcion	string		`json:"descripcion"`
	Enlinea     bool		`json:"enlinea"`
	Url         string		`json:"url"`
	BlobKey		appengine.BlobKey	`json:"imgurl"`
}


func init() {
    http.HandleFunc("/wsdetalle", ShowOfDetalle)
}

func ShowOfDetalle(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Add(time.Duration(-18000)*time.Second)
	var timetolive = 900 //seconds
	c := appengine.NewContext(r)
	var b []byte
	var d detalle
	if item, err := memcache.Get(c, "d_"+r.FormValue("id")); err == memcache.ErrCacheMiss {
		oferta, _ := model.GetOferta(c, r.FormValue("id"))
		if now.After(oferta.FechaHoraPub) {
			d.IdEmp = oferta.IdEmp
			d.IdOft = oferta.IdOft
			d.IdCat = oferta.IdCat
			d.Oferta = oferta.Oferta
			d.Empresa = oferta.Empresa
			d.Descripcion = oferta.Descripcion
			d.Enlinea = oferta.Enlinea
			d.Url = oferta.Url
			d.BlobKey = oferta.BlobKey
		}

		b, _ = json.Marshal(d)
		item := &memcache.Item{
			Key:   "d_"+r.FormValue("id"),
			Value: b,
			Expiration: time.Duration(timetolive)*time.Second,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("Memcache.Add d_idoft : %v", err)
		}
	} else {
		//c.Infof("memcache retrieve d_idoft : %v", r.FormValue("id"))
		b = item.Value
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

