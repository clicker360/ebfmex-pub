package oferta

import (
    "appengine"
	"encoding/json"
	"sortutil"
    "net/http"
	"model"
	"time"
)

type WsSucursal struct{
	IdOft       string `json:"idoft"`
	IdEmp       string `json:"idemp"`
	IdSuc       string `json:"idsuc"`
	Sucursal    string `json:"sucursal"`
	FechaHora   time.Time `json:"timestamp"`
	Status		string `json:"status"`
}

func init() {
    //http.HandleFunc("/delofsuc", DelOfSuc)
    http.HandleFunc("/r/ofsuc", ShowEmpSucursalOft)
}

func DelOfSuc(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsSucursal
	out.IdSuc = r.FormValue("idsuc")
	out.IdOft = r.FormValue("idoft")
	err := model.DelOfertaSucursal(c, out.IdOft, out.IdSuc)
	if err != nil {
		out.Status = "notFound"
	} else {
		out.Status = "ok"
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}

func ShowOfSucursales(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	ofsucs, _ := model.GetOfertaSucursales(c, r.FormValue("id"))
	wssucs := make([]WsSucursal, 0 ,len(*ofsucs))
	for i,v:= range *ofsucs {
		wssucs[i].IdOft = v.IdOft
		wssucs[i].IdSuc = v.IdSuc
		wssucs[i].IdEmp = v.IdEmp
		wssucs[i].Sucursal = v.Sucursal
		wssucs[i].FechaHora = v.FechaHora
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(wssucs)
	w.Write(b)
}

/*
	Listado de sucursales por empresa
*/
func ShowEmpSucs(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	emsucs := model.GetEmpSucursales(c, r.FormValue("IdEmp"))
	if emsucs != nil {
		wssucs := make([]WsSucursal, len(*emsucs), cap(*emsucs))
		for i,v:= range *emsucs {
			wssucs[i].IdOft = ""
			wssucs[i].IdSuc = v.IdSuc
			wssucs[i].IdEmp = v.IdEmp
			wssucs[i].Sucursal = v.Nombre
			wssucs[i].FechaHora = v.FechaHora
		}
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(wssucs)
		w.Write(b)
	}
}

/*
	Listado de sucursales por empresa con la oferta marcada
*/
func ShowEmpSucursalOft(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	emsucs := model.GetEmpSucursales(c, r.FormValue("idemp"))
	if emsucs != nil {
		ofsucs, _ := model.GetOfertaSucursales(c, r.FormValue("idoft"))
		wssucs := make([]WsSucursal, len(*emsucs), cap(*emsucs))
		for i,es:= range *emsucs {
			for _,os:= range *ofsucs {
				if os.IdSuc == es.IdSuc {
					wssucs[i].IdOft = os.IdOft
				}
			}
			wssucs[i].IdSuc = es.IdSuc
			wssucs[i].IdEmp = es.IdEmp
			wssucs[i].Sucursal = es.Nombre
			wssucs[i].FechaHora = es.FechaHora
		}
		sortutil.AscByField(wssucs, "Sucursal")

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(wssucs)
		w.Write(b)
	}
}
