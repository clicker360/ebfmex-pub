package oferta

import (
	"appengine"
	"appengine/datastore"
	"strings"
	"model"
	"time"
)

func init() {
}

func putSearchData(c appengine.Context, value string, key *datastore.Key, idoft string, idcat int, enlinea bool) {
	r := strings.NewReplacer("."," ",","," ",";"," ",":"," ","!"," ","~"," ","¿"," ","?"," ","#"," ","_"," ","-"," ","+"," ","'"," ","\""," ","*"," ","$"," ","("," ",")"," ","="," ","%"," ","&"," ","<"," ",">"," ","|"," ","@"," ","·"," ","["," ","]"," ","{"," ","}"," ","¡"," ","!"," ", "\n", " ", "\r", " ", "\t", " ")
	if err := model.DelOfertaSearchData(c, key); err != nil {
		c.Errorf("Datastore Delete Kind:SearchData, key:%s", key)
	} else {
		for _, v := range strings.Split(r.Replace(value), " ") {
			if(len(v)>3) {
				w := strings.ToLower(v)
				if(model.ValidSearchData.MatchString(w)) {
					sd := &model.SearchData {
						Sid: key.Encode(),
						Kind: "Oferta",
						Field: "Descripcion",
						Value: w,
						IdCat: idcat,
						Enlinea: enlinea,
						FechaHora: time.Now(),
					}
					// Pa llave es el idoft + palabra clave
					_, err := datastore.Put(c, datastore.NewKey(c, "SearchData", idoft+w, 0, nil), sd)
					if err != nil {
						c.Errorf("Datastore Put Kind:SearchData, IdOft:%s, word:%s", idoft, w)
					}
				} else {
					c.Infof("Intento de palabra inválida al diccionario, word:%s", idoft, w)
				}
			}
		}
	}
}

