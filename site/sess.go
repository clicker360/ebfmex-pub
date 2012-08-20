package site

import (
    "appengine"
    "appengine/datastore"
    "appengine/mail"
    "net/http"
	"crypto/md5"
	"strings"
	"bytes"
    "time"
	"fmt"
    "io"
)

func init() {
    http.HandleFunc("/acceso", Acceso)
    http.HandleFunc("/recupera", Recover)
    http.HandleFunc("/salir", Salir)
}

type Sess struct {
	Md5			string
	Id			string
	User		string
	Name		string
	Expiration	time.Time
	ErrMsg		string `datastore:"-"`
	ErrClass	string `datastore:"-"`
}

func setSess(w http.ResponseWriter, c appengine.Context, key *datastore.Key, email string, name string) (string, *datastore.Key, error) {
	h := md5.New()
	io.WriteString(h, key.Encode())
	io.WriteString(h, fmt.Sprintf("%s", time.Now()))
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	ex := time.Now().AddDate(0,0,5)
	s := Sess{
		Md5:		md5,
		Id:			key.Encode(),
		User:		email,
		Name:		name,
		Expiration:	ex,
	}
	cKey, err := datastore.Put(c, datastore.NewKey(c, "Sess", email, 0, nil), &s)
	if err != nil {
		return "", nil, err
	}
	// Se crean 2 cookies, una con el key de sesión otra con el número random llave 
	//csc := http.Cookie{ Name: "ebfmex-pub-sesscontrol-ua", Value: md5, Expires: ex, Path: "/" }
	//http.SetCookie(w, &csc)
	//csc = http.Cookie{ Name: "ebfmex-pub-sessid-ua", Value: cKey.Encode(), Expires: ex, Path: "/" }
	//http.SetCookie(w, &csc)
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sesscontrol-ua=%s; expires=%s; path=/;", md5, ex.Format("Wed, 07-Oct-2012 14:23:42 GMT")))
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sessid-ua=%s; expires=%s; path=/;", cKey.Encode(), ex.Format("Wed, 07-Oct-2012 14:23:42 GMT")))
	return md5, cKey, err
}

func IsSess(w http.ResponseWriter, r *http.Request, c appengine.Context) (Sess, bool) {
	var s Sess
	if ck, err := r.Cookie("ebfmex-pub-sessid-ua"); err == nil {
		key, _ := datastore.DecodeKey(ck.Value)
		if err := datastore.Get(c, key, &s); err != nil {
			// no hay sesión
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			/* se verifica el control de sesion */
			if cr, err := r.Cookie("ebfmex-pub-sesscontrol-ua"); err == nil {
				if s.Md5 != cr.Value {
					// Md5 no coincide, intenta entrar con otra cookie
					return s, false
				} else if time.Now().After(s.Expiration) {
					// Sesión expirada
					return s, false
				}
				// OK
				// Hay sesión
				return s, true
			}
		}
	}
	return s, false
}

func Acceso(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	var st Sess
	if _, ok := IsSess(w, r, c); !ok {
		//fmt.Fprintf(w, "u:%s, p:%s", r.FormValue("u"), r.FormValue("p"))
		if(r.FormValue("u") != "" && r.FormValue("p") != "") {
			/* validar usuario y pass */
			if(validEmail.MatchString(r.FormValue("u")) && validName.MatchString(r.FormValue("p"))) {
				q := datastore.NewQuery("Cta").Filter("Email =", r.FormValue("u")).Filter("Pass =", r.FormValue("p")).Filter("Status =", true)
				if count, _ := q.Count(c); count != 0 {
					for t := q.Run(c); ; {
						var g Cta
						key, err := t.Next(&g)
						if err == datastore.Done {
							break
						}
						// Coincide contraseña, se activa una sesión
						_, _, err = setSess(w, c, key, g.Email, g.Nombre)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						// Redireccion
						http.Redirect(w, r, "/dash", http.StatusFound)
						return
					}
				}
			}
			st.User = r.FormValue("u")
			st.ErrMsg = "Usuario y/o Contraseña no aceptados"
			st.ErrClass = "show"
		} else {
			st.User = r.FormValue("u")
			st.ErrMsg = "Proporcione usuario y contraseña"
			st.ErrClass = "show"
		}
	} else {
		// hay sesión
		http.Redirect(w, r, "/dash", http.StatusFound)
		return
	}
	tc := make(map[string]interface{})
	tc["Sess"] = st
    accesoErrorTpl.ExecuteTemplate(w, "cta", tc)
}

func Salir(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := IsSess(w, r, c); ok {
		s.Expiration = time.Now().AddDate(-1,0,0)
		_, err := datastore.Put(c, datastore.NewKey(c, "Sess", s.User, 0, nil), &s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sesscontrol-ua=%s; expires=%s; path=/;", "", "Wed, 07-Oct-2000 14:23:42 GMT"))
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sessid-ua=%s; expires=%s; path=/;", "", "Wed, 07-Oct-2000 14:23:42 GMT"))
	http.Redirect(w, r, "/registro", http.StatusFound)
}

func Recover(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := IsSess(w, r, c); !ok {
		var email string = strings.TrimSpace(r.FormValue("Email"))
		var rfc string = strings.TrimSpace(r.FormValue("RFC"))
		if email != "" && validEmail.MatchString(email) && rfc != "" && validRfc.MatchString(rfc) {
			// intenta buscar en la base un usuario con email y empresa
			if cta, err := GetCta(c, email); err == nil {
				q := datastore.NewQuery("Empresa").Filter("RFC =", rfc).Ancestor(cta.Key(c)).Limit(3)
				if count, _ := q.Count(c); count != 0 {
					for t := q.Run(c); ; {
						_, err := t.Next(&cta)
						if err == datastore.Done {
							break
						}
						var hbody bytes.Buffer
						var sender string
						if (appengine.AppID(c) == "ebfmxorg") {
							sender =  "El Buen Fin <contacto@elbuenfin.org>"
						} else {
							sender =  "El Buen Fin <ahuezo@clicker360.com>"
						}
						if err := mailRecoverTpl.Execute(&hbody, cta); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						// Coincide email y RFC, se manda correo con contraseña
						msg := &mail.Message{
							Sender:		sender,
							To:			[]string{cta.Email},
							Subject:	"Recuperación de contraseña / El Buen Fin",
							HTMLBody:	hbody.String(),
						}
						if err := mail.Send(c, msg); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						} else {
							http.Redirect(w, r, "/recoverok.html", http.StatusFound)
							return
						}
						//fmt.Fprintf(w, mailRecover, cta.Email, cta.Pass)
						return
					}
				}
			}
		}
		http.Redirect(w, r, "/nocta.html", http.StatusFound)
		return
	} else {
		http.Redirect(w, r, "/dash", http.StatusFound)
	}
}

