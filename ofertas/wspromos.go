package oferta

import (
    "appengine"
	"encoding/json"
	"sortutil"
    "net/http"
	"model"
)

type Promo struct{
	Token	string `json:"token"`
	Id		string `json:"id"`
	Tipo	string `json:"tipo"`
	Status  string `json:"status"`
}

func init() {
    //http.HandleFunc("/modpromo", ModPromo)
    //http.HandleFunc("/delpromo", ModPromo)
    //http.HandleFunc("/promosxo", ShowPromo)
    //http.HandleFunc("/tcxo", ShowTcxo)
    //http.HandleFunc("/modtcxo", ModTcxo)
}

func ModPromo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out Promo
	out.Token = r.FormValue("token")
	out.Tipo = r.FormValue("tipo")
	out.Id = r.FormValue("id")
	oferta := model.GetOferta(c, out.Id)
	if oferta.IdEmp != "none" {
		if out.Tipo == "desc" {
			oferta.Descuento = out.Token
		} else if out.Tipo == "promo" {
			oferta.Promocion = out.Token
		} else if out.Tipo == "meses" {
			oferta.Meses = out.Token
		}
		err := model.PutOferta(c, oferta)
		if err != nil {
			out.Status = "writeErr"
		} else {
			out.Status = "ok"
		}
	} else {
		out.Status = "notFound"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

func ShowPromo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	oferta := model.GetOferta(c, r.FormValue("id"))
	var out [3]Promo
	if oferta.Descuento != "" {
		out[0].Token = oferta.Descuento
		out[0].Id = oferta.IdOft
		out[0].Tipo = "desc"
	}
	if oferta.Promocion != "" {
		out[1].Token = oferta.Promocion
		out[1].Id = oferta.IdOft
		out[1].Tipo = "promo"
	}
	if oferta.Meses != "" {
		out[2].Token = oferta.Meses
		out[2].Id = oferta.IdOft
		out[2].Tipo = "meses"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

type jsonTc struct {
	Id			int		`json: "id"`
	Tarjeta		string	`json: "tarjeta"`
	Selected	int		`json: "selected"`
	Status		string	`json: "status"`
}

func ShowTcxo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	oferta := model.GetOferta(c, r.FormValue("id"))
	var b []byte
	if oferta.IdEmp != "none" {
		var tarjetas []jsonTc
		if err := json.Unmarshal(oferta.Tarjetas, &tarjetas); err != nil {
			model.Check(err)
		}
		sortutil.AscByField(tarjetas, "Tarjeta")
		b, _ = json.Marshal(tarjetas)
	} else {
		b = []byte("[]")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}


/*
	Las tarjetas se procesan como un simple string json
	se guardan como tal y se leen como tal
	Es responsabilidad del cliente interpretarlas
	Si se agregan tarjetas al catálogo no debería importarle al 
	datastore
	Posteriormente se puede hacer un proceso batch para indexarlas
	en una tabla independiente.
*/
func ModTcxo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var jsonBlob = []byte(r.FormValue("token"))
	var out jsonTc
	/* 
		el arreglo tarjetas sólo se usa para validar el json de entrada
	*/
	var tarjetas []jsonTc
	if err := json.Unmarshal(jsonBlob, &tarjetas); err != nil {
		out.Status = "syntaxErr"
	} else {
		oferta := model.GetOferta(c, r.FormValue("id"))
		if oferta.IdEmp != "none" {
			oferta.Tarjetas = []byte(r.FormValue("token"))
			err := model.PutOferta(c, oferta)
			if err != nil {
					out.Status = "writeErr"
			} else {
				out.Status = "ok"
			}
		} else {
			out.Status = "notFound"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}


