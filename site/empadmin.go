package site

import (
    "appengine"
    "appengine/datastore"
	"strings"
    "net/http"
	"strconv"
    "time"
//	"fmt"
)

type FormDataEmp struct {
	IdEmp			string
	RFC				string
	ErrRFC			string
	Nombre			string
	ErrNombre		string
	RazonSoc		string
	ErrRazonSoc		string
	DirCalle		string
	ErrDirCalle		string
	DirCol			string
	ErrDirCol		string
	DirEnt			string
	ErrDirEnt		string
	Entidades		*[]Entidad
	Organismos		*[]Organismo
	DirMun			string
	ErrDirMun		string
	DirCp			string
	ErrDirCp		string
	NumSuc			string
	ErrNumSuc		string
	OrgEmp			string
	ErrOrgEmp		string
	OrgEmpOtro		string
	ErrOrgEmpOtro	string
	OrgEmpReg		string
	ErrOrgEmpReg	string
	Url				string
	ErrUrl			string
	PartLinea		int
	ExpComer		int
	Desc			string
	ErrDesc			string
}

func init() {
    http.HandleFunc("/se", GetEmp)
    http.HandleFunc("/ne", NewEmp)
    http.HandleFunc("/me", ModEmp)
    http.HandleFunc("/de", DelEmp)
    http.HandleFunc("/le", ShowListEmp)
}

func ShowListEmp(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		u, _ := GetCta(c, s.User)
		tc := make(map[string]interface{})
		tc["Sess"] = s
		tc["Empresa"] = listEmp(c, u)
		empadmTpl.ExecuteTemplate(w, "empresa", tc)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func GetEmp(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		u, _ := GetCta(c, s.User)
		e, err := u.GetEmpresa(c, r.FormValue("IdEmp"))
		if err == datastore.ErrNoSuchEntity {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//Se trata de una empresa nueva
		}
		formEmp(c, w, &s, e)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func ModEmp(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		u, _ := GetCta(c, s.User)
		// formato con validación 
		_, valid := formatoEmp(w, r, s, true)
		if !valid {return}

		fe := fillEmpresa(r)
		_, err := u.PutEmpresa(c, &fe)
		if err == datastore.ErrNoSuchEntity {
			// Aviso NO EXISTE EMPRESA
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		ShowListEmp(w, r)
	} else {
		defer http.Redirect(w, r, "/", http.StatusFound)
	}
}

func NewEmp(w http.ResponseWriter, r *http.Request) {
	// formato con validación 
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		_, valid := formatoEmp(w, r, s, true)
		if !valid { return }
		u, _ := GetCta(c, s.User)
		fe := fillEmpresa(r)

		// Se añade una empresa
		//e, err := u.AddEmpresa(c, &fe)
		_, err := u.AddEmpresa(c, &fe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//formEmp(c, w, &s, e)
		ShowListEmp(w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func DelEmp(w http.ResponseWriter, r *http.Request) {
	// formato con validación 
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		u, _ := GetCta(c, s.User)
		if _, err := u.DelEmpresa(c, r.FormValue("IdEmp")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ShowListEmp(w, r)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func formEmp(c appengine.Context, w http.ResponseWriter, s *Sess, e *Empresa) {
	fd := empresaToForm(*e)
	fd.Entidades = listEnt(c, e.DirEnt)
	fd.Organismos = listOrg(c, e.OrgEmp)
	tc := make(map[string]interface{})
	tc["Sess"] = s
	tc["FormDataEmp"] = fd
	empadmTpl.ExecuteTemplate(w, "empresa", tc)
	return
}

func listEmp(c appengine.Context, u *Cta) *[]Empresa {
	q := datastore.NewQuery("Empresa").Ancestor(u.Key(c)).Limit(50)
	empresas := make([]Empresa, 0, 50)
	if _, err := q.GetAll(c, &empresas); err != nil {
		return nil
	}
	return &empresas
}

func listEnt(c appengine.Context, ent string) *[]Entidad {
	q := datastore.NewQuery("Entidad")
	estados := make([]Entidad, 0, 32)
	if _, err := q.GetAll(c, &estados); err != nil {
		return nil
	}
	for i, _ := range estados {
		if(ent == estados[i].CveEnt) {
			estados[i].Selected = `selected`
		}
	}
	return &estados
}

func listOrg(c appengine.Context, siglas string) *[]Organismo {
	q := datastore.NewQuery("Organismo")
	orgs := make([]Organismo, 0, 32)
	if _, err := q.GetAll(c, &orgs); err != nil {
		return nil
	}
	for i, _ := range orgs {
		if(siglas == orgs[i].Siglas) {
			orgs[i].Selected = `selected`
		}
	}
	return &orgs
}

func formatoEmp(w http.ResponseWriter, r *http.Request, s Sess, valida bool) (FormDataEmp, bool){
	c := appengine.NewContext(r)
	partlinea, _ := strconv.Atoi(r.FormValue("PartLinea"))
	expcomer, _ := strconv.Atoi(r.FormValue("ExpComer"))
	fd := FormDataEmp {
		IdEmp: strings.TrimSpace(r.FormValue("IdEmp")),
		RFC: strings.TrimSpace(r.FormValue("RFC")),
		ErrRFC: "",
		Nombre: strings.TrimSpace(r.FormValue("Nombre")),
		ErrNombre: "",
		RazonSoc: strings.TrimSpace(r.FormValue("RazonSoc")),
		ErrRazonSoc: "",
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
		NumSuc: strings.TrimSpace(r.FormValue("NumSuc")),
		ErrNumSuc: "",
		OrgEmp: strings.TrimSpace(r.FormValue("OrgEmp")),
		ErrOrgEmp: "",
		OrgEmpOtro: strings.TrimSpace(r.FormValue("OrgEmpOtro")),
		ErrOrgEmpOtro: "",
		OrgEmpReg: strings.TrimSpace(r.FormValue("OrgEmpReg")),
		ErrOrgEmpReg: "",
		Url	: strings.TrimSpace(r.FormValue("Url")),
		ErrUrl:	"",
		PartLinea: partlinea,
		ExpComer: expcomer,
		Desc: strings.TrimSpace(r.FormValue("Desc")),
	}
	if valida {
		var ef bool
		ef = false
		if fd.RFC == "" || !validRfc.MatchString(fd.RFC) {
			fd.ErrRFC = "invalid"
			ef = true
		}
		if fd.Nombre == "" || !validName.MatchString(fd.Nombre) {
			fd.ErrNombre = "invalid"
			ef = true
		}
		if fd.RazonSoc == "" || !validName.MatchString(fd.RazonSoc) {
			fd.ErrRazonSoc = "invalid"
			ef = true
		}
		/*
		if fd.DirEnt ==  || !validName.MatchString(fd.DirEnt) {
			fd.ErrDirEnt = "invalid"
			ef = true
		}
		*/
		if fd.DirMun == "" || !validSimpleText.MatchString(fd.DirMun) {
			fd.ErrDirMun = "invalid"
			ef = true
		}
		if fd.DirCalle == "" || !validSimpleText.MatchString(fd.DirCalle) {
			fd.ErrDirCalle = "invalid"
			ef = true
		}
		if fd.DirCol == "" || !validSimpleText.MatchString(fd.DirCol) {
			fd.ErrDirCol = "invalid"
			ef = true
		}
		if fd.DirCp == "" || !validCP.MatchString(fd.DirCp) {
			fd.ErrDirCp = "invalid"
			ef = true
		}
		if fd.NumSuc == "" || !validNum.MatchString(fd.NumSuc) {
			fd.ErrNumSuc = "invalid"
			ef = true
		}
		if fd.OrgEmp != "" && !validSimpleText.MatchString(fd.OrgEmp) {
			fd.ErrOrgEmp = "invalid"
			ef = true
		}
		if fd.OrgEmpOtro != "" && !validSimpleText.MatchString(fd.OrgEmpOtro) {
			fd.ErrOrgEmpOtro = "invalid"
			ef = true
		}
		if fd.OrgEmpReg != "" && !validSimpleText.MatchString(fd.OrgEmpReg) {
			fd.ErrOrgEmpReg = "invalid"
			ef = true
		}
		if fd.Url != "" && !validUrl.MatchString(fd.Url) {
			fd.ErrUrl = "invalid"
			ef = true
		}
		/*
		if fd.Desc != "" && !validSimpleText.MatchString(fd.Desc) {
			fd.ErrDesc = "invalid"
			ef = true
		}
		*/

		if ef {
			fd.Entidades = listEnt(c, strings.TrimSpace(r.FormValue("DirEnt")))
			fd.Organismos = listOrg(c, fd.OrgEmp)
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormDataEmp"] = fd
			empadmTpl.ExecuteTemplate(w, "empresa", tc)
			return fd, false
		}
	}
	return fd, true
}

func fillEmpresa(r *http.Request) Empresa {
	partlinea, _ := strconv.Atoi(r.FormValue("PartLinea"))
	expcomer, _ := strconv.Atoi(r.FormValue("ExpComer"))
	e := Empresa {
		IdEmp:		strings.TrimSpace(r.FormValue("IdEmp")),
		RFC:		strings.TrimSpace(r.FormValue("RFC")),
		Nombre:		strings.TrimSpace(r.FormValue("Nombre")),
		RazonSoc:	strings.TrimSpace(r.FormValue("RazonSoc")),
		DirCalle:	strings.TrimSpace(r.FormValue("DirCalle")),
		DirCol:		strings.TrimSpace(r.FormValue("DirCol")),
		DirEnt:		strings.TrimSpace(r.FormValue("DirEnt")),
		DirMun:		strings.TrimSpace(r.FormValue("DirMun")),
		DirCp:		strings.TrimSpace(r.FormValue("DirCp")),
		NumSuc:		strings.TrimSpace(r.FormValue("NumSuc")),
		OrgEmp:		strings.TrimSpace(r.FormValue("OrgEmp")),
		OrgEmpOtro:	strings.TrimSpace(r.FormValue("OrgEmpOtro")),
		OrgEmpReg:	strings.TrimSpace(r.FormValue("OrgEmpReg")),
		Url:		strings.TrimSpace(r.FormValue("Url")),
		PartLinea:  partlinea,
		ExpComer:	expcomer,
		Desc:		strings.TrimSpace(r.FormValue("Desc")),
		FechaHora:	time.Now(),
		Status:		true,
	}
	return e
}

func empresaToForm(e Empresa) FormDataEmp {
	fe := FormDataEmp {
		IdEmp:		e.IdEmp,
		RFC:		e.RFC,
		Nombre:		e.Nombre,
		RazonSoc:	e.RazonSoc,
		DirCalle:	e.DirCalle,
		DirCol:		e.DirCol,
		DirEnt:		e.DirEnt,
		DirMun:		e.DirMun,
		DirCp:		e.DirCp,
		NumSuc:		e.NumSuc,
		OrgEmp:		e.OrgEmp,
		OrgEmpOtro:	e.OrgEmpOtro,
		OrgEmpReg:	e.OrgEmpReg,
		PartLinea:  e.PartLinea,
		ExpComer:	e.ExpComer,
		Desc:		e.Desc,
		Url:		e.Url,
	}
	return fe
}
