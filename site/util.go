package site

import (
		"html/template"
        "math/rand"
        "regexp"
        "time"
)

// randId returns a string of random letters.
func randId(l int) string {
        n := make([]byte, l)
        for i := range n {
                n[i] = 'a' + byte(rand.Intn(26))
        }
        return string(n)
}

func init() {
        // Seed number generator with the current time.
        rand.Seed(time.Now().UnixNano())
}



// validName matches a valid name string.
var validNum = regexp.MustCompile(`^[0-9]+$`)
var validCP = regexp.MustCompile(`^[0-9]{5,5}`)
var validName = regexp.MustCompile(`^[a-zA-Z áéíóúAÉÍÓÚÑñäëïöü\.\'\-]+$`)
var validSimpleText = regexp.MustCompile(`^[a-zA-Z0-9].+$`)
var validPass = regexp.MustCompile(`^[a-zA-Z0-9 áéíóúAÉÍÓÚÑñäëïöü¡¿\.\,\;\:\!\{\}\~\(\)\?\#\_\+\/\%\$\'\"\*\-]+$`)
var validEmail = regexp.MustCompile(`^([0-9a-zA-Z]([-.\w]*[0-9a-zA-Z])*@(([0-9a-zA-Z])+([-\w]*[0-9a-zA-Z])*\.)+[a-zA-Z]{2,9})$`)
var validTel = regexp.MustCompile(`^([\(]{1}[0-9]{2,3}[\)]{1}[\.| |\-]{0,1}|^[0-9]{3,4}[\.|\-| ]?)?[0-9]{3,4}(\.|\-| )?[0-9]{3,4}$`)
//var validRfc = regexp.MustCompile(`^[A-Z,Ñ,&]{3,4}[0-9]{2}[0-1][0-9][0-3][0-9][A-Z,0-9]?[A-Z,0-9]?[0-9,A-Z]?`)
var validRfc = regexp.MustCompile(`^([A-Z&Ññ]{3}|[A-Z][AEIOU][A-Z]{2})\d{2}((01|03|05|07|08|10|12)(0[1-9]|[12]\d|3[01])|02(0[1-9]|[12]\d)|(04|06|09|11)(0[1-9]|[12]\d|30))([A-Z0-9]{2}[0-9A])?$`)
var validUrl = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \?=.-]*)*\/?$`)
const IdLen = 8;

var registroTpl = template.Must(template.ParseFiles("templates/registro.html")) //, "templates/login.html"))
var registroErrorTpl = template.Must(template.ParseFiles("templates/registro_aviso.html"))
var accesoErrorTpl = template.Must(template.ParseFiles("templates/acceso_error.html"))
var ctadmTpl = template.Must(template.ParseFiles("templates/ctadm.html"))
var dashTpl = template.Must(template.ParseFiles("templates/dashboard.html"))
var activaTpl = template.Must(template.ParseFiles("templates/activacion.html"))
var listUsersTpl = template.Must(template.ParseFiles("templates/list_users.html"))
var listSessTpl = template.Must(template.ParseFiles("templates/list_sess.html"))
var mailActivationCodeTpl = template.Must(template.ParseFiles("templates/activation_code.html"))
var mailAvisoActivacionTpl = template.Must(template.ParseFiles("templates/mail_aviso_activacion.html"))
var activationMessageTpl = template.Must(template.ParseFiles("templates/activation_result.html"))
var mailRecoverTpl = template.Must(template.ParseFiles("templates/mail_recover.html"))
var empadmTpl = template.Must(template.ParseFiles("templates/empadm.html"))
var proxTpl = template.Must(template.ParseFiles("templates/proximamente.html"))
