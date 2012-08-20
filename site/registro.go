package site

import (
    "appengine"
    "appengine/datastore"
	"appengine/mail"
    "appengine/user"
    "net/http"
	"html/template"
	"crypto/md5"
	"bytes"
    "time"
	"fmt"
	"io"
	"sharded_counter"
	"model"
	"sess"
)

type urlCfm struct {
	Md5			string
	Nombre		string
	Apellidos	string
	Email		string
	FechaHora	time.Time
	Llave		string
	AppId		string
}

func init() {
    http.HandleFunc("/registrar", Registrar)
    http.HandleFunc("/c", ConfirmaCodigo)
}

func Registrar(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	s, ok := sess.IsSess(w, r, c)
	if ok {
		http.Redirect(w, r, "/cta", http.StatusFound)
		return
	}
	fd, valid := ctaForm(w, r, s, true, registroTpl)
	if valid {
		u, err := model.GetCta(c, fd.Email)
		ctaFill(r, u)
		if err != nil {
			// No hay Cuenta registrada
			u.FechaHora = time.Now()
			u.Status = false

			// Generar código de confirmación distindo cada vez. Md5 del email + fecha-hora
			h := md5.New()
			io.WriteString(h, fmt.Sprintf("%s%s%s%s", time.Now(), u.Email, u.Pass, model.RandId(12)))
			u.CodigoCfm = fmt.Sprintf("%x", h.Sum(nil))
		}

		// Se agrega la cuenta sin activar para realizar el proceso de código de confirmación
		if u, err = model.PutCta(c, u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//Si hay estatus es que ya existe
		if(u.Status == false) {
			/* No se ha activado, por tanto se inicia el proceso de código de verificación */
			m := urlCfm{
				Md5:		u.CodigoCfm,
				Nombre:		u.Nombre,
				Email:		u.Email,
				FechaHora:	time.Now(),
				Llave:		u.Key(c).Encode(),
				AppId:		appengine.AppID(c),
			}
			var hbody bytes.Buffer
			var sender string
			if (appengine.AppID(c) == "ebfmxorg") {
				sender =  "El Buen Fin <contacto@elbuenfin.org>"
			} else {
				sender =  "El Buen Fin <ahuezo@clicker360.com>"
			}
			// Envia código activación 
			if err := mailActivationCodeTpl.Execute(&hbody, m); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			msg := &mail.Message{
				Sender:  sender,
				To:      []string{m.Email},
				Subject: "Codigo de Activación de Registro / El Buen Fin en línea",
				HTMLBody: hbody.String(),
			}
			if err := mail.Send(c, msg); err != nil {
				/* Problemas para enviar el correo NOK */
				http.Error(w, err.Error(), http.StatusInternalServerError)
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				if err := activationMessageTpl.ExecuteTemplate(w, "codesend", m); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

			// ************************************************************
			// Si el hay usuario admin se despliega el código de activación
			// ************************************************************
			if gu := user.Current(c); gu != nil {
				if err := mailActivationCodeTpl.Execute(w, m); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			return
		}
		if err := registroErrorTpl.Execute(w, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ConfirmaCodigo(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	md5 := r.FormValue("m")
    key, _ := datastore.DecodeKey(r.FormValue("c"))
	var g model.Cta
    if err := datastore.Get(c, key, &g); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	/* Se verifica el código de confirmación */
	if(g.CodigoCfm == md5 && g.Status == false) {
		// Si se confirma el md5 la cuenta admin se le asigna un folio y se activa el status
		if err := sharded_counter.Increment(c, "cuenta_admin"); err == nil {
			if folio, err := sharded_counter.Count(c, "cuenta_admin"); err == nil {
				g.Folio = folio
				g.Status = true
				g.CodigoCfm = "Confirmado"
				_, err := datastore.Put(c, key, &g)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				/* Prende la sesion */
				_, _, err = sess.SetSess(w, c, key, g.Email, g.Nombre)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Envia código activación 
				var hbody bytes.Buffer
				var sender string
				if (appengine.AppID(c) == "ebfmxorg") {
					sender =  "El Buen Fin <contacto@elbuenfin.org>"
				} else {
					sender =  "El Buen Fin <ahuezo@clicker360.com>"
				}
				if err := mailAvisoActivacionTpl.Execute(&hbody, g); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				msg := &mail.Message{
						Sender:  sender,
						To:      []string{g.Email},
						Subject: "Cuenta Activada / El Buen Fin en línea",
						HTMLBody:	hbody.String(),
				}
				if err := mail.Send(c, msg); err != nil {
					// Problemas para enviar el correo NOK 
					http.Error(w, err.Error(), http.StatusInternalServerError)
					http.Redirect(w, r, "/", http.StatusFound)
				}
				// avisa del éxito independientemente del correo
				if err := activationMessageTpl.ExecuteTemplate(w, "confirm", g); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			} else {
				// El Folio no es seguro, se deshecha la operación o se encola
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// El Folio no es seguro, se deshecha la operación o se encola 
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := activationMessageTpl.ExecuteTemplate(w, "codeerr", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

var registroTpl = template.Must(template.ParseFiles("templates/registro.html")) //, "templates/login.html"))
var registroErrorTpl = template.Must(template.ParseFiles("templates/registro_aviso.html"))
var mailActivationCodeTpl = template.Must(template.ParseFiles("templates/activation_code.html"))
var activationMessageTpl = template.Must(template.ParseFiles("templates/activation_result.html"))
var mailAvisoActivacionTpl = template.Must(template.ParseFiles("templates/mail_aviso_activacion.html"))
