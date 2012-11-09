package backend

import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "model"
	//"strings"
	"strconv"
)


func init() {
    http.HandleFunc("/backend/updtctaemp", updateCtaEmpresa)
}

func updateCtaEmpresa(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	const batch = 200

	page,_ := strconv.Atoi(r.FormValue("pg"))
	if page < 1 {
		page = 1
	}
	offset := batch * (page - 1)
    q := datastore.NewQuery("Cta").Offset(offset).Limit(batch)
	n,_ := q.Count(c)
	for cursor := q.Run(c); ; {
		var cta model.Cta
		_, err := cursor.Next(&cta)
		if err == datastore.Done  {
			break
		}

		q2 := datastore.NewQuery("Empresa").Ancestor(cta.Key(c))
		for cursor := q2.Run(c); ; {
			var emp model.Empresa
			_, err := cursor.Next(&emp)
			if err == datastore.Done  {
				break
			}
			var ce model.CtaEmpresa
			ce.IdEmp = emp.IdEmp
			ce.Email = cta.Email
			ce.EmailAlt = cta.EmailAlt
			//_, err1 := datastore.Put(c, datastore.NewKey(c, "CtaEmpresa", ce.Email+"_"+ce.IdEmp, 0, nil), &ce)
			_, err1 := datastore.Put(c, datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil), &ce)
			if err1 != nil {
				c.Errorf("PutCtaEmpresa(); Error al intentar actualizar CtaEmpresa : %v", emp.IdEmp)
			}
		}
	}
	c.Infof("UpdateServingLogoUrl() Pagina: %d, actualizados: %d, del %d al %d", page, n, offset, offset+n)
}

