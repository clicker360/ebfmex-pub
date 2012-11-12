package backend

import (
	"appengine"
    "appengine/datastore"
	"net/http"
	"strconv"
	"strings"
	"model"
)

func init() {
    http.HandleFunc("/backend/updtsrvlogo", UpdateServingUrlLogo)
}

func UpdateServingUrlLogo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	const batch = 200
	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
	q := datastore.NewQuery("Oferta").Offset(offset).Order("-FechaHora").Limit(batch)
	n,_ := q.Count(c)
	for i := q.Run(c); ; {
		var e model.Oferta
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}

		//Oferta.Promocion es URL s180
		//Oferta.Descuento es URL s70
		e.Promocion = ""
		e.Descuento = ""
		emplogo := model.GetLogo(c, e.IdEmp)
		if emplogo != nil {
			// Tenga lo que tenga, se pasa Sp4 a Oferta.Promocion
			if(emplogo.Sp4 != "")  {
				e.Promocion = emplogo.Sp4
				e.Descuento = strings.Replace(emplogo.Sp4, "s180", "s70",1)
			}
			//c.Errorf("Get Cta Key; Error al intentar leer key.Parent() de Empresa : %v", e.IdEmp)
			_, err = datastore.Put(c, key, &e)
			if err != nil {
				c.Errorf("PutOferta(); Error al intentar actualizar oferta : %v", e.IdOft)
			}
		}

	}
	c.Infof("UpdateServingLogo() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
	return
}
