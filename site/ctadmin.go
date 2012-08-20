package site

import (
	"html/template"
    "appengine"
    "appengine/datastore"
    "net/http"
	"strings"
    "time"
	"fmt"
)

type FormDataCta struct {
	Folio			string `datastore:"-"`
	Nombre			string `datastore:"-"`
	ErrNombre		string `datastore:"-"`
	Apellidos		string `datastore:"-"`
	ErrApellidos	string `datastore:"-"`
	Puesto			string `datastore:"-"`
	ErrPuesto		string `datastore:"-"`
	Email			string `datastore:"-"`
	ErrEmail		string `datastore:"-"`
	EmailAlt		string `datastore:"-"`
	ErrEmailAlt		string `datastore:"-"`
	Pass			string `datastore:"-"`
	ErrPass			string `datastore:"-"`
	Pass1			string `datastore:"-"`
	ErrPass1		string `datastore:"-"`
	Tel				string `datastore:"-"`
	ErrTel			string `datastore:"-"`
	Cel				string `datastore:"-"`
	ErrCel			string `datastore:"-"`
	TermCond		string `datastore:"-"`
	ErrTermCond		string `datastore:"-"`
}

func init() {
    http.HandleFunc("/registro", registro)
    //http.HandleFunc("/registro", mantenimiento)
    http.HandleFunc("/dash", dash)
    http.HandleFunc("/cta", CtaShow)
    http.HandleFunc("/ctamod", CtaMod)
    http.HandleFunc("/ctadel", CtaDel)
}

func mantenimiento(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404.html", http.StatusFound)
}

func registro(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := IsSess(w, r, c); !ok {
		var fd FormDataCta
		tc := make(map[string]interface{})
		//tc["Sess"] = s
		tc["FormDataCta"] = fd
		registroTpl.ExecuteTemplate(w, "cta", tc)
		return
	} else {
		http.Redirect(w, r, "/dash", http.StatusFound)
	}
}

func dash(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		if g, err := GetCta(c, s.User); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormDataCta"] = ctaToForm(*g)
			dashTpl.ExecuteTemplate(w, "cta", tc)
			return
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func CtaShow(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		if g, err := GetCta(c, s.User); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormDataCta"] = ctaToForm(*g)
			ctadmTpl.ExecuteTemplate(w, "cta", tc)
			return
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func CtaMod(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		if u, err := GetCta(c, s.User); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if fd, valid := ctaForm(w, r, s, true, ctadmTpl); !valid { 
				return 
			} else {
				ctaFill(r, u)
				if _, err := PutCta(c, u); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				tc := make(map[string]interface{})
				tc["Sess"] = s
				tc["FormDataCta"] = fd
				ctadmTpl.ExecuteTemplate(w, "cta", tc)
			}
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func CtaDel(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		if u, err := GetCta(c, s.User); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			/* Desactiva Status */
			if(r.FormValue("desactiva")=="1") {
				s.Expiration = time.Now().AddDate(-1,0,0)
				_, err := datastore.Put(c, datastore.NewKey(c, "Sess", s.User, 0, nil), &s)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				u.CodigoCfm = ""
				u.Status = false
				_, err = datastore.Put(c, u.Key(c), &u)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				tc := make(map[string]interface{})
				tc["CanceledCta"] = 1
				ctadmTpl.ExecuteTemplate(w, "cta", tc)
				w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sesscontrol-ua=%s; expires=%s; path=/;", "", "Wed, 07-Oct-2000 14:23:42 GMT"))
				w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sessid-ua=%s; expires=%s; path=/;", "", "Wed, 07-Oct-2000 14:23:42 GMT"))
				return
			}
			tc := make(map[string]interface{})
			tc["AskCancelCta"] = 1
			ctadmTpl.ExecuteTemplate(w, "cta", tc)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func ctaForm(w http.ResponseWriter, r *http.Request, s Sess, valida bool, tpl *template.Template) (FormDataCta, bool){
	fd := FormDataCta {
		Folio: strings.TrimSpace(r.FormValue("IdEmp")),
		Nombre: strings.TrimSpace(r.FormValue("Nombre")),
		ErrNombre: "",
		Apellidos: strings.TrimSpace(r.FormValue("Apellidos")),
		ErrApellidos: "",
		Puesto: strings.TrimSpace(r.FormValue("Puesto")),
		ErrPuesto: "",
		Email: strings.TrimSpace(r.FormValue("Email")),
		ErrEmail: "",
		EmailAlt: strings.TrimSpace(r.FormValue("EmailAlt")),
		ErrEmailAlt: "",
		Pass: strings.TrimSpace(r.FormValue("Pass")),
		ErrPass: "",
		Pass1: strings.TrimSpace(r.FormValue("Pass1")),
		ErrPass1: "",
		Tel: strings.TrimSpace(r.FormValue("Tel")),
		ErrTel: "",
		Cel: strings.TrimSpace(r.FormValue("Cel")),
		ErrCel: "",
		TermCond: strings.TrimSpace(r.FormValue("TermCond")),
		ErrTermCond: "",
	}
	if valida {
		var ef bool
		ef = false
		if fd.Nombre == "" || !validName.MatchString(fd.Nombre) {
			fd.ErrNombre = "invalid"
			ef = true
		}
		if fd.Apellidos == "" || !validName.MatchString(fd.Apellidos) {
			fd.ErrApellidos = "invalid"
			ef = true
		}
		if fd.Puesto != "" && !validSimpleText.MatchString(fd.Puesto) {
			fd.ErrPuesto = "invalid"
			ef = true
		}
		if fd.Email == "" || !validEmail.MatchString(fd.Email) {
			fd.ErrEmail = "invalid"
			ef = true
		}
		if fd.EmailAlt != "" && !validEmail.MatchString(fd.EmailAlt) {
			fd.ErrEmailAlt = "invalid"
			ef = true
		}
		if r.FormValue("modificar") != "1" {
			if (fd.Pass != fd.Pass1 || fd.Pass == "" || fd.Pass1 == "" || !validPass.MatchString(fd.Pass)) {
				fd.ErrPass = "invalid"
				fd.ErrPass1 = "invalid"
				ef = true
			}
		} else {
			if ((fd.Pass != fd.Pass1 || !validPass.MatchString(fd.Pass)) && (fd.Pass != "" || fd.Pass1 != "")) {
			//if (fd.Pass != fd.Pass1 && fd.Pass != "" && fd.Pass1 != "" && !validPass.MatchString(fd.Pass)) {
				fd.ErrPass = "invalid"
				fd.ErrPass1 = "invalid"
				ef = true
			}
		}
		if fd.Tel == "" || !validTel.MatchString(fd.Tel) {
			fd.ErrTel = "invalid"
			ef = true
		}
		if fd.Cel != "" && !validTel.MatchString(fd.Cel) {
			fd.ErrCel = "invalid"
			ef = true
		}
		if fd.TermCond == "" && r.FormValue("t") == "r" {
			fd.ErrTermCond = "invalid"
			ef = true
		}

		if ef {
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormDataCta"] = fd
			tpl.ExecuteTemplate(w, "cta", tc)
			return fd, false
		}
	}
	return fd, true
}

func ctaFill(r *http.Request, cta *Cta) {
	cta.Nombre=		strings.TrimSpace(r.FormValue("Nombre"))
	cta.Apellidos=	strings.TrimSpace(r.FormValue("Apellidos"))
	cta.Puesto=		strings.TrimSpace(r.FormValue("Puesto"))
	cta.Email=		strings.TrimSpace(r.FormValue("Email"))
	cta.EmailAlt=	strings.TrimSpace(r.FormValue("EmailAlt"))
	if r.FormValue("Pass") != "" {
		cta.Pass=		strings.TrimSpace(r.FormValue("Pass"))
	}
	cta.Tel=		strings.TrimSpace(r.FormValue("Tel"))
	cta.Cel=		strings.TrimSpace(r.FormValue("Cel"))
}

func ctaToForm(e Cta) FormDataCta {
	fd := FormDataCta {
		Nombre:		e.Nombre,
		Apellidos:	e.Apellidos,
		Puesto:		e.Puesto,
		Email:		e.Email,
		EmailAlt:	e.EmailAlt,
		Pass:		e.Pass,
		Tel:		e.Tel,
		Cel:		e.Cel,
	}
	return fd
}
