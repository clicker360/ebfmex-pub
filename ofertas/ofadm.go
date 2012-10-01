// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	"appengine/datastore"
	//"resize"
	//"bytes"
	"strings"
	"strconv"
	//"fmt"
	//"image"
	//"image/jpeg"
	_ "image/png" // import so we can read PNG files.
	//"io"
	"net/http"
	"html/template"
	"sess"
	"model"
	"time"
)

type FormDataOf struct {
	IdOft			string
	IdEmp			string
	IdCat			int
	Categorias		*[]model.Categoria
	Empresa			string
	Oferta			string
	ErrOferta		string
	Descripcion		string
	ErrDescripcion	string
	Codigo			string
	ErrCodigo		string
	Precio			string
	ErrPrecio		string
	Descuento		string
	ErrDescuento	string
	Enlinea			bool
	Url				string
	ErrUrl			string
	Tarjetas		string // Texto separado por comas
	ErrTarjetas		string
	Meses			string
	ErrMeses		string
	FechaHoraPub    time.Time
	ErrFechaHoraPub string
	StatusPub		bool
	FechaHora		time.Time
	Ackn			string
}

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {
	http.HandleFunc("/of", model.ErrorHandler(OfShow))
	http.HandleFunc("/ofs", model.ErrorHandler(OfShowList))
	//http.HandleFunc("/ofimgup", model.ErrorHandler(ImgOferta))
	http.HandleFunc("/ofmod", model.ErrorHandler(OfMod))
	http.HandleFunc("/ofdel", model.ErrorHandler(OfDel))
	//http.HandleFunc("/ofsucadd", model.ErrorHandler(OfAddSucursal))
	//http.HandleFunc("/ofsucdel", model.ErrorHandler(OfDelSucursal))
}

func OfShowList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		//usuario, _ := model.GetCta(c, s.User)
		tc := make(map[string]interface{})
		tc["Sess"] = s
		if empresa := model.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
			tc["Empresa"] = empresa
			tc["Oferta"] = listOf(c, empresa.IdEmp)
		}
		ofadmTpl.ExecuteTemplate(w, "ofertas", tc)
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

func listOf(c appengine.Context, IdEmp string) *[]model.Oferta {
	q := datastore.NewQuery("Oferta").Filter("IdEmp =", IdEmp).Limit(500)
	n, _ := q.Count(c)
	ofertas := make([]model.Oferta, 0, n)
	if _, err := q.GetAll(c, &ofertas); err != nil {
		return nil
	}
	return &ofertas
}

func listCat(c appengine.Context, IdCat int) *[]model.Categoria {
	q := datastore.NewQuery("Categoria")
	n, _ := q.Count(c)
	categorias := make([]model.Categoria, 0, n)
	if _, err := q.GetAll(c, &categorias); err != nil {
		return nil
	}
	for i, _ := range categorias {
		if(IdCat == categorias[i].IdCat) {
			categorias[i].Selected = `selected`
		}
	}
	return &categorias
}

func OfShow(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		oferta := model.GetOferta(c, r.FormValue("IdOft"))
		var id string
		if oferta.IdEmp != "none" {
			id = oferta.IdEmp
		} else {
			id = r.FormValue("IdEmp")
		}
		fd := ofToForm(*oferta)
		if empresa := model.GetEmpresa(c, id); empresa != nil {
			tc["Empresa"] = empresa
			fd.IdEmp = empresa.IdEmp
			oferta.Empresa = empresa.Nombre
		}
		fd.Categorias = listCat(c, oferta.IdCat);
		tc["FormDataOf"] = fd
		ofadmTpl.ExecuteTemplate(w, "oferta", tc)
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

// Modifica si hay, Crea si no hay
// Requiere IdEmp. IdOft es opcional, si no hay lo crea, si hay modifica
func OfMod(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		tc := make(map[string]interface{})
		tc["Sess"] = s
		fd, valid :=ofForm(w, r, true)
		ofertamod := oftFill(fd)
		oferta := model.GetOferta(c, ofertamod.IdOft)
		if empresa := model.GetEmpresa(c, ofertamod.IdEmp); empresa != nil {
			tc["Empresa"] = empresa
			fd.IdEmp = empresa.IdEmp
			fd.Empresa = empresa.Nombre
			ofertamod.Empresa = empresa.Nombre
		}
		// TODO
		// es preferible poner un regreso avisando que no existe la empresa
		if valid {
			if oferta.IdOft != "none" {
				// Ya existe
				err := model.PutOferta(c, &ofertamod)
				model.Check(err)
			} else {
				// nueva oferta
				err := model.NewOferta(c, &ofertamod)
				model.Check(err)
			}
			fd = ofToForm(ofertamod)
			fd.Ackn = "Ok";
		}
		fd.Categorias = listCat(c, ofertamod.IdCat);
		tc["FormDataOf"] = fd
		ofadmTpl.ExecuteTemplate(w, "oferta", tc)
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

func OfDel(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := sess.IsSess(w, r, c); ok {
		if err := model.DelOferta(c, r.FormValue("IdOft")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		OfShowList(w, r)
		return
	}
	http.Redirect(w, r, "/ofertas", http.StatusFound)
}

func ofForm(w http.ResponseWriter, r *http.Request, valida bool) (FormDataOf, bool){
	c := appengine.NewContext(r)
	var fh time.Time
	if r.FormValue("FechaHoraPub") != "" {
		fh, _ = time.Parse("2012-08-21 14:01:41", strings.TrimSpace(r.FormValue("FechaHoraPub")))
	} else {
		fh = time.Now()
	}
	el, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("Enlinea")))
	st, _ := strconv.ParseBool(strings.TrimSpace(r.FormValue("StatusPub")))
	ic, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("IdCat")))
	fd := FormDataOf {
		IdOft:			strings.TrimSpace(r.FormValue("IdOft")),
		IdEmp:			strings.TrimSpace(r.FormValue("IdEmp")),
		IdCat:			ic,
		Oferta:			strings.TrimSpace(r.FormValue("Oferta")),
		ErrOferta: "",
		Descripcion:	strings.TrimSpace(r.FormValue("Descripcion")),
		ErrDescripcion: "",
		Codigo:			strings.TrimSpace(r.FormValue("Codigo")),
		ErrCodigo: "",
		Precio:			strings.TrimSpace(r.FormValue("Precio")),
		ErrPrecio: "",
		Descuento:		strings.TrimSpace(r.FormValue("Descuento")),
		ErrDescuento: "",
		Enlinea:		el,
		Url:			strings.TrimSpace(r.FormValue("Url")),
		ErrUrl: "",
		Tarjetas:		strings.TrimSpace(r.FormValue("Tarjetas")),
		ErrTarjetas: "",
		Meses:			strings.TrimSpace(r.FormValue("Meses")),
		ErrMeses: 	"",
		FechaHoraPub:	fh,
		ErrFechaHoraPub: "",
		StatusPub:		st,
	}
	if valida {
		var ef bool
		ef = false
		if fd.Oferta == "" || !model.ValidSimpleText.MatchString(fd.Oferta) {
			fd.ErrOferta = "invalid"
			ef = true
		}
		if fd.Descripcion != "" && !model.ValidSimpleText.MatchString(fd.Descripcion) && len(fd.Descripcion) > 200 {
			fd.ErrDescripcion = "invalid"
			ef = true
		}
		if fd.Precio != "" && !model.ValidPrice.MatchString(fd.Precio) {
			fd.ErrPrecio = "invalid"
			ef = true
		}
		if fd.Descuento != "" && !model.ValidPercent.MatchString(fd.Descuento) {
			fd.ErrDescuento = "invalid"
			ef = true
		}
		if fd.Url != "" && !model.ValidUrl.MatchString(fd.Url) {
			fd.ErrUrl = "invalid"
			ef = true
		}

		fd.Categorias = listCat(c, ic);
		if ef {
			return fd, false
		}
	}
	return fd, true
}

func oftFill(fd FormDataOf) model.Oferta {
	s := model.Oferta {
		IdOft:			fd.IdOft,
		IdEmp:			fd.IdEmp,
		IdCat:			fd.IdCat,
		Oferta:			fd.Oferta,
		Descripcion:	fd.Descripcion,
		Codigo:			fd.Codigo,
		Precio:			fd.Precio,
		Descuento:		fd.Descuento,
		Enlinea:		fd.Enlinea,
		Url:			fd.Url,
		Tarjetas:		fd.Tarjetas,
		Meses:			fd.Meses,
		FechaHoraPub:	fd.FechaHoraPub,
		StatusPub:		fd.StatusPub,
		FechaHora:		time.Now(),
	}
	return s;
}

func ofToForm(e model.Oferta) FormDataOf {
	fd := FormDataOf {
		IdOft:			e.IdOft,
		IdEmp:			e.IdEmp,
		IdCat:			e.IdCat,
		Oferta:			e.Oferta,
		Descripcion:	e.Descripcion,
		Codigo:			e.Codigo,
		Precio:			e.Precio,
		Descuento:		e.Descuento,
		Enlinea:		e.Enlinea,
		Url:			e.Url,
		Tarjetas:		e.Tarjetas,
		Meses:			e.Meses,
		FechaHoraPub:	e.FechaHoraPub,
		StatusPub:		e.StatusPub,
	}
	return fd
}

var ofadmTpl = template.Must(template.ParseFiles("templates/ofadm.html"))
