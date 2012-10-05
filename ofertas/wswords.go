package oferta

import (
    "appengine"
	"encoding/json"
    "net/http"
	"model"
)

type Word struct{
	Token	string `json:"token"`
	Id		string `json:"id"`
	Status  string `json:"status"`
}

func init() {
    http.HandleFunc("/addword", AddWord)
    http.HandleFunc("/delword", DelWord)
    http.HandleFunc("/wordsxo", ShowWords)
}

func AddWord(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out Word
	out.Token = r.FormValue("token")
	out.Id = r.FormValue("id")
	oferta := model.GetOferta(c, out.Id)
	if oferta.IdEmp != "none" {
		var palabra model.OfertaPalabra
		palabra.IdOft = out.Id
		palabra.Palabra = out.Token
		err := oferta.PutOfertaPalabra(c, &palabra)
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

func DelWord(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out Word
	out.Token = r.FormValue("token")
	out.Id = r.FormValue("id")
	oferta := model.GetOferta(c, out.Id)
	if oferta.IdEmp != "none" {
		err := model.DelOfertaPalabra(c, out.Id)
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

func ShowWords(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ofertas, _ := model.GetOfertaPalabras(c, r.FormValue("id"))
	words := make([]Word, 0 ,len(*ofertas))
	for i,v:= range *ofertas {
		words[i].Id = v.IdOft
		words[i].Token = v.Palabra
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(words)
	w.Write(b)
}
