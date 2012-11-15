package sucursales

import (
	"html/template"
    "appengine"
    "appengine/datastore"
    "net/http"
	"sortutil"
	"strings"
	//"fmt"
	"time"
	"model"
	"sess"
)

type FormDataSuc struct {
	IdSuc			string
	IdEmp			string
	Nombre			string
	ErrNombre		string
	Tel				string
	ErrTel			string
	DirCalle		string
	ErrDirCalle		string
	DirCol			string
	ErrDirCol		string
	DirEnt			string
	ErrDirEnt		string
	Entidades		*[]model.Entidad
	DirMun			string
	ErrDirMun		string
	DirCp			string
	ErrDirCp		string
	GeoUrl			string
	Geo1			string
	Geo2			string
	Geo3			string
	Geo4			string
	Ackn			string
}

func init() {
    http.HandleFunc("/r/sucursales", ShowListSuc)
    http.HandleFunc("/r/sucursal", SucShow)
    http.HandleFunc("/r/sucmod", SucMod)
    http.HandleFunc("/r/sucdel", SucDel)
}

func ShowListSuc(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		tc := make(map[string]interface{})
		tc["Sess"] = s
		empresa, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tc["Empresa"] = empresa
		tc["Sucursal"] = listSuc(c, u, empresa.IdEmp)
		sucadmTpl.ExecuteTemplate(w, "sucursales", tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func listSuc(c appengine.Context, u *model.Cta, IdEmp string) *[]model.Sucursal {
	q := datastore.NewQuery("Sucursal").Ancestor(datastore.NewKey(c, "Empresa", IdEmp, 0, u.Key(c)))
	n, _ := q.Count(c)
	sucursales := make([]model.Sucursal, 0, n)
	if _, err := q.GetAll(c, &sucursales); err != nil {
		return nil
	}
	sortutil.AscByField(sucursales, "Nombre")
	return &sucursales
}

func SucShow(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		tc := make(map[string]interface{})
		tc["Sess"] = s
		sucursal := model.GetSuc(c, u, r.FormValue("IdSuc"), r.FormValue("IdEmp"))
		var id string
		if sucursal.IdEmp != "none" {
			id = sucursal.IdEmp
		} else {
			id = r.FormValue("IdEmp")
		}
		empresa, err := u.GetEmpresa(c, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tc["Empresa"] = empresa
		fd := sucToForm(*sucursal)
		fd.Entidades = model.ListEnt(c, sucursal.DirEnt)
		tc["FormDataSuc"] = fd
		sucadmTpl.ExecuteTemplate(w, "sucursal", tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

// Modifica si hay, Crea si no hay
// Requiere IdEmp. IdSuc es opcional, si no hay lo crea, si hay modifica
func SucMod(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		u, _ := model.GetCta(c, s.User)
		tc := make(map[string]interface{})
		tc["Sess"] = s
		fd, valid :=sucForm(w, r, true)
		empresa, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sucursal := sucFill(r)
		if valid {
			if empresa != nil {
				newsuc, err := empresa.PutSuc(c, &sucursal)
				//fmt.Fprintf(w, "IdSuc: %s", newsuc.IdSuc);
				fd = sucToForm(*newsuc)
				fd.Ackn = "Ok";
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		fd.Entidades = model.ListEnt(c, strings.TrimSpace(r.FormValue("DirEnt")))
		tc["Empresa"] = empresa
		tc["FormDataSuc"] = fd
		sucadmTpl.ExecuteTemplate(w, "sucursal", tc)
	} else {
		http.Redirect(w, r, "/r/registro", http.StatusFound)
	}
}

func SucDel(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := sess.IsSess(w, r, c); ok {
		if err := model.DelSuc(c, r.FormValue("IdSuc")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ShowListSuc(w, r)
		return
	}
	http.Redirect(w, r, "/r/sucursales", http.StatusFound)
}

func sucForm(w http.ResponseWriter, r *http.Request, valida bool) (FormDataSuc, bool){
	fd := FormDataSuc {
		IdSuc: strings.TrimSpace(r.FormValue("IdSuc")),
		IdEmp: strings.TrimSpace(r.FormValue("IdEmp")),
		GeoUrl: strings.TrimSpace(r.FormValue("GeoUrl")),
		Geo1: strings.TrimSpace(r.FormValue("Geo1")),
		Geo2: strings.TrimSpace(r.FormValue("Geo2")),
		Geo3: strings.TrimSpace(r.FormValue("Geo3")),
		Geo4: strings.TrimSpace(r.FormValue("Geo4")),
		Nombre: strings.TrimSpace(r.FormValue("Nombre")),
		ErrNombre: "",
		Tel: strings.TrimSpace(r.FormValue("Tel")),
		ErrTel: "",
		DirCalle: strings.TrimSpace(r.FormValue("DirCalle")),
		ErrDirCalle: "",
		DirCol: strings.TrimSpace(r.FormValue("DirCol")),
		ErrDirCol: "",
		DirEnt: strings.TrimSpace(r.FormValue("DirEnt")),
		ErrDirEnt: "",
		DirMun: strings.TrimSpace(r.FormValue("DirMun")),
		ErrDirMun: "",
		DirCp: strings.TrimSpace(r.FormValue("DirCp")),
		ErrDirCp: "",
	}
	if valida {
		var ef bool
		ef = false
		if fd.Nombre == "" || !model.ValidSimpleText.MatchString(fd.Nombre) {
			fd.ErrNombre = "invalid"
			ef = true
		}
		if fd.Tel != "" && !model.ValidTel.MatchString(fd.Tel) {
			fd.ErrTel = "invalid"
			ef = true
		}
		if fd.DirEnt == "" || !model.ValidSimpleText.MatchString(fd.DirEnt) {
			fd.ErrDirEnt = "invalid"
			ef = true
		}
		if fd.DirMun == "" || !model.ValidSimpleText.MatchString(fd.DirMun) {
			fd.ErrDirMun = "invalid"
			ef = true
		}
		if fd.DirCalle == "" || !model.ValidSimpleText.MatchString(fd.DirCalle) {
			fd.ErrDirCalle = "invalid"
			ef = true
		}
		if fd.DirCol == "" || !model.ValidSimpleText.MatchString(fd.DirCol) {
			fd.ErrDirCol = "invalid"
			ef = true
		}
		/*
		if fd.DirCp == "" || !model.ValidCP.MatchString(fd.DirCp) {
			fd.ErrDirCp = "invalid"
			ef = true
		}
		*/

		if ef {
			return fd, false
		}
	}
	return fd, true
}

func sucFill(r *http.Request) model.Sucursal {
	s := model.Sucursal {
		IdEmp:		strings.TrimSpace(r.FormValue("IdEmp")),
		IdSuc:		strings.TrimSpace(r.FormValue("IdSuc")),
		Nombre:		strings.TrimSpace(r.FormValue("Nombre")),
		Tel:		strings.TrimSpace(r.FormValue("Tel")),
		DirCalle:	strings.TrimSpace(r.FormValue("DirCalle")),
		DirCol:		strings.TrimSpace(r.FormValue("DirCol")),
		DirEnt:		strings.TrimSpace(r.FormValue("DirEnt")),
		DirMun:		strings.TrimSpace(r.FormValue("DirMun")),
		DirCp:		strings.TrimSpace(r.FormValue("DirCp")),
		GeoUrl:		strings.TrimSpace(r.FormValue("GeoUrl")),
		Geo1:		strings.TrimSpace(r.FormValue("Geo1")),
		Geo2:		strings.TrimSpace(r.FormValue("Geo2")),
		Geo3:		strings.TrimSpace(r.FormValue("Geo3")),
		Geo4:		strings.TrimSpace(r.FormValue("Geo4")),
		FechaHora:	time.Now().Add(time.Duration(model.GMTADJ)*time.Second),
	}
	return s;
}

func sucToForm(e model.Sucursal) FormDataSuc {
	fd := FormDataSuc {
		IdSuc:		e.IdSuc,
		IdEmp:		e.IdEmp,
		Nombre:		e.Nombre,
		Tel:		e.Tel,
		DirCalle:	e.DirCalle,
		DirCol:		e.DirCol,
		DirEnt:		e.DirEnt,
		DirMun:		e.DirMun,
		DirCp:		e.DirCp,
		GeoUrl:		e.GeoUrl,
		Geo1:		e.Geo1,
		Geo2:		e.Geo2,
		Geo3:		e.Geo3,
		Geo4:		e.Geo4,
	}
	return fd
}

var sucadmTpl = template.Must(template.ParseFiles("templates/sucadm.html"))
