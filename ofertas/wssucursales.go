package oferta

import (
    "appengine"
	"encoding/json"
    "net/http"
	"strconv"
	"model"
	"time"
)

type WsSucursal struct{
	IdOft       string `json:"idoft"`
	IdEmp       string `json:"idemp"`
	IdSuc       string `json:"idsuc"`
	Sucursal    string `json:"sucursal"`
	Status		string `json:"status"`
}

func init() {
    http.HandleFunc("/addofsuc", AddOfSuc)
    http.HandleFunc("/delofsuc", DelOfSuc)
    http.HandleFunc("/ofsuc", ShowEmpSucursalOft)
}

func AddOfSuc(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out WsSucursal
	out.IdSuc = r.FormValue("idsuc")
	out.IdOft = r.FormValue("idoft")
	out.IdEmp = r.FormValue("idemp")
	oferta := model.GetOferta(c, out.IdOft)
	suc := model.GetSuc(c, out.IdSuc)
	if oferta.IdEmp != "none" && suc.IdEmp != "none" {
		lat, _ := strconv.ParseFloat(suc.Geo1, 64)
		lng, _ := strconv.ParseFloat(suc.Geo2, 64)
		var ofsuc model.OfertaSucursal
		out.Sucursal = suc.Nombre
		ofsuc.IdOft = out.IdOft
		ofsuc.IdSuc = out.IdSuc
		ofsuc.IdEmp = out.IdEmp
		ofsuc.Sucursal = suc.Nombre
		ofsuc.Lat = lat
		ofsuc.Lng = lng
		ofsuc.Empresa = oferta.Empresa
		ofsuc.Oferta = oferta.Oferta
		ofsuc.Descripcion = oferta.Descripcion
		ofsuc.Promocion = oferta.Promocion
		ofsuc.Precio = oferta.Precio
		ofsuc.Descuento = oferta.Descuento
		ofsuc.Url = oferta.Url
		ofsuc.StatusPub = oferta.StatusPub
		ofsuc.FechaHora = time.Now()

		err := oferta.PutOfertaSucursal(c, &ofsuc)
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
	sucs := make([]WsSucursal, 0 ,len(*ofsucs))
	for i,v:= range *ofsucs {
		sucs[i].IdOft = v.IdOft
		sucs[i].IdSuc = v.IdSuc
		sucs[i].IdEmp = v.IdEmp
		sucs[i].Sucursal = v.Sucursal
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(sucs)
	w.Write(b)
}

/*
	Listado de sucursales por empresa
*/
func ShowEmpSucs(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	emsucs := model.GetEmpSucursales(c, r.FormValue("IdEmp"))
	wssucs := make([]WsSucursal, len(*emsucs), cap(*emsucs))
	for i,v:= range *emsucs {
		wssucs[i].IdOft = ""
		wssucs[i].IdSuc = v.IdSuc
		wssucs[i].IdEmp = v.IdEmp
		wssucs[i].Sucursal = v.Nombre
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(wssucs)
	w.Write(b)
}

/*
	Listado de sucursales por empresa con la oferta marcada
*/
func ShowEmpSucursalOft(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	emsucs := model.GetEmpSucursales(c, r.FormValue("idemp"))
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
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(wssucs)
	w.Write(b)
}
