// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	"appengine/datastore"
	"appengine/blobstore"
	"sortutil"
	"strings"
	"strconv"
	_ "image/png" // import so we can read PNG files.
	"net/http"
	"net/url"
	"html/template"
	"sess"
	"model"
	"time"
	"io"
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
	Tarjetas		[]byte // json
	ErrTarjetas		string
	Meses			string
	ErrMeses		string
	FechaHoraPub    time.Time
	ErrFechaHoraPub string
	StatusPub		bool
	FechaHora		time.Time
	Ackn			string
	Sucursales		string // cadena de id's de sucursales separadas por espacio
	UploadURL		*url.URL
	BlobKey			appengine.BlobKey
}

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.

func init() {
	http.HandleFunc("/of", model.ErrorHandler(OfShow))
	http.HandleFunc("/ofs", model.ErrorHandler(OfShowList))
	http.HandleFunc("/ofmod", model.ErrorHandler(OfMod))
	http.HandleFunc("/ofdel", model.ErrorHandler(OfDel))
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
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
	sortutil.AscByField(ofertas, "Oferta")
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

		/*
		 * Se crea el form para el upload del blob
		 */
		uploadURL, err := blobstore.UploadURL(c, "/ofimgup", nil)
		if err != nil {
			serveError(c, w, err)
			return
		}
		fd.UploadURL = uploadURL

		tc["FormDataOf"] = fd
		w.Header().Set("Content-Type", "text/html")
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
		var fd FormDataOf
		var valid bool
		var ofertamod model.Oferta

		if  r.FormValue("IdOft") == "new" {
			if empresa := model.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
				tc["Empresa"] = empresa
				fd.IdEmp = empresa.IdEmp
				fd.Empresa = empresa.Nombre
				ofertamod.IdEmp = empresa.IdEmp
				ofertamod.Oferta = "Nueva oferta"
				ofertamod.FechaHora = time.Now().Add(time.Duration(-18000)*time.Second) // 5 horas menos
				ofertamod.FechaHoraPub = time.Now().Add(time.Duration(-18000)*time.Second) // 5 horas menos
				ofertamod.Empresa = strings.ToUpper(empresa.Nombre)
				ofertamod.BlobKey = "none"
				o, err := model.NewOferta(c, &ofertamod)
				model.Check(err)
				fd = ofToForm(*o)
			fd.Ackn = "Ok";
			} else {
				// redireccionar
				http.Redirect(w, r, "/le?d=o", http.StatusFound)
			}
		} else {
			/* 
			 * Se pide un id oferta que en teoría existe, se consulta y se cambia
			 * Se valida y si no existe se informa un error
			 */
			fd, valid =ofForm(w, r, true)

			ofertamod.IdOft = fd.IdOft
			ofertamod.IdEmp = fd.IdEmp
			ofertamod.IdCat = fd.IdCat
			ofertamod.Oferta = fd.Oferta
			ofertamod.Descripcion = fd.Descripcion
			ofertamod.Enlinea =	fd.Enlinea
			ofertamod.Url =	fd.Url
			ofertamod.FechaHoraPub = fd.FechaHoraPub
			ofertamod.StatusPub = fd.StatusPub
			//ofertamod.BlobKey = fd.BlobKey
			ofertamod.FechaHora = time.Now().Add(time.Duration(-18000)*time.Second)

			oferta := model.GetOferta(c, ofertamod.IdOft)
			if oferta.IdOft != "none" {
				if empresa := model.GetEmpresa(c, ofertamod.IdEmp); empresa != nil {
					tc["Empresa"] = empresa
					fd.IdEmp = empresa.IdEmp
					fd.Empresa = empresa.Nombre
					ofertamod.Empresa = strings.ToUpper(empresa.Nombre)
					ofertamod.BlobKey = oferta.BlobKey
				}
				// TODO
				// es preferible poner un regreso avisando que no existe la empresa
				if valid {
					// Ya existe
					err := model.PutOferta(c, &ofertamod)
					model.Check(err)

					// Se borran las relaciones oferta-sucursal
					err = model.DelOfertaSucursales(c, oferta.IdOft)
					model.Check(err)

					// Se reconstruyen las Relaciones oferta-sucursal con las solicitadas
					idsucs := strings.Fields(r.FormValue("schain"))
					for _, idsuc := range idsucs {
						suc := model.GetSuc(c, idsuc)

						lat, _ := strconv.ParseFloat(suc.Geo1, 64)
						lng, _ := strconv.ParseFloat(suc.Geo2, 64)

						var ofsuc model.OfertaSucursal
						ofsuc.IdOft = ofertamod.IdOft
						ofsuc.IdSuc = idsuc
						ofsuc.IdEmp = ofertamod.IdEmp
						ofsuc.Sucursal = suc.Nombre
						ofsuc.Lat = lat
						ofsuc.Lng = lng
						ofsuc.Empresa = ofertamod.Empresa
						ofsuc.Oferta = ofertamod.Oferta
						ofsuc.Descripcion = ofertamod.Descripcion
						ofsuc.Promocion = ofertamod.Promocion
						ofsuc.Precio = ofertamod.Precio
						ofsuc.Descuento = ofertamod.Descuento
						ofsuc.Url = ofertamod.Url
						ofsuc.StatusPub = ofertamod.StatusPub
						ofsuc.FechaHora = time.Now().Add(time.Duration(-18000)*time.Second)

						err := ofertamod.PutOfertaSucursal(c, &ofsuc)
						model.Check(err)

					}
					fd = ofToForm(ofertamod)
					fd.Ackn = "Ok";
				}
			} else {
				// no existe la oferta
			}
		}

		fd.Categorias = listCat(c, ofertamod.IdCat);

		/*
		 * Se crea el form para el upload del blob
		 */
		uploadURL, err := blobstore.UploadURL(c, "/ofimgup", nil)
		if err != nil {
			serveError(c, w, err)
			return
		}
		fd.UploadURL = uploadURL
		tc["FormDataOf"] = fd

		w.Header().Set("Content-Type", "text/html")
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
		fh, _ = time.Parse("_2 Jan 15:04:05", strings.TrimSpace(r.FormValue("FechaHoraPub"))+" 00:00:00")
		fh = fh.AddDate(2012,0,0)
	} else {
		fh = time.Now().Add(time.Duration(-18000)*time.Second) // 5 horas menos
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
		//Precio:			strings.TrimSpace(r.FormValue("Precio")),
		//ErrPrecio: "",
		//Descuento:		strings.TrimSpace(r.FormValue("Descuento")),
		ErrDescuento: "",
		Enlinea:		el,
		Url:			strings.TrimSpace(r.FormValue("Url")),
		ErrUrl: "",
		//Tarjetas:		strings.TrimSpace(r.FormValue("Tarjetas")),
		//ErrTarjetas: "",
		//Meses:			strings.TrimSpace(r.FormValue("Meses")),
		//Promocion:		strings.TrimSpace(r.FormValue("Promocion")),
		ErrMeses:	"",
		FechaHoraPub:	fh,
		ErrFechaHoraPub: strings.TrimSpace(fh.Format("_2 Jan")),
		StatusPub:		st,
	}
	if valida {
		var ef bool
		ef = false
		if fd.Oferta == "" || !model.ValidSimpleText.MatchString(fd.Oferta) {
			fd.ErrOferta = "invalid"
			ef = true
		}
		if fd.Descripcion == "" || !model.ValidSimpleText.MatchString(fd.Descripcion) && len(fd.Descripcion) > 200 {
			fd.ErrDescripcion = "invalid"
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

func ofToForm(e model.Oferta) FormDataOf {
	fd := FormDataOf {
		IdOft:			e.IdOft,
		IdEmp:			e.IdEmp,
		IdCat:			e.IdCat,
		Oferta:			e.Oferta,
		Descripcion:	e.Descripcion,
		//Codigo:			e.Codigo,
		//Precio:			e.Precio,
		//Descuento:		e.Descuento,
		//Promocion:		e.Promocion,
		Enlinea:		e.Enlinea,
		Url:			e.Url,
		//Tarjetas:		e.Tarjetas,
		//Meses:			e.Meses,
		FechaHoraPub:	e.FechaHoraPub,
		ErrFechaHoraPub:	strings.TrimSpace(e.FechaHoraPub.Format("_2 Jan")),
		StatusPub:		e.StatusPub,
		BlobKey:		e.BlobKey,
	}
	return fd
}

var ofadmTpl = template.Must(template.ParseFiles("templates/ofadm.html"))
