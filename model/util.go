package model

import (
        "math/rand"
		"net/http"
        "time"
		"regexp"
		"html/template"
)

// check aborts the current execution if err is non-nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// errorHandler wraps the argument handler with an error-catcher that
// returns a 500 HTTP error if the request fails (calls check with err non-nil).
func ErrorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if _, ok := recover().(error); ok {
				w.WriteHeader(http.StatusInternalServerError)
				tc := make(map[string]interface{})
				tc["ErrMsg"] = "Error de ejecución"
				ErrorGeneralTpl.Execute(w, tc)
			}
		}()
		fn(w, r)
	}
}
var ErrorGeneralTpl = template.Must(template.ParseFiles("templates/aviso_error_general.html"))

// randId returns a string of random letters.
func RandId(l int) string {
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
var ValidNum = regexp.MustCompile(`^[0-9]+$`)
var ValidCP = regexp.MustCompile(`^[0-9]{5,5}`)
var ValidKey = regexp.MustCompile(`^[a-zA-Z]+$`)
var ValidName = regexp.MustCompile(`^[a-zA-Z áéíóúAÉÍÓÚÑñäëïöü\.\'\-]+$`)
var ValidSimpleText = regexp.MustCompile(`^[a-zA-Z0-9 _áéíóúAÉÍÓÚÑñäëïöü¡¿ªº\&\.\,\;\:\!\{\}\~\(\)\?\#\_\+\/\%\$\'\"\*\-\<\>\[\]\@\|]+$`)
//var ValidSimpleText = regexp.MustCompile(`^[a-zA-Z0-9].+$`)
var ValidPass = regexp.MustCompile(`^[a-zA-Z0-9 áéíóúAÉÍÓÚÑñäëïöü¡¿\.\,\;\:\!\{\}\~\(\)\?\#\_\+\/\%\$\'\"\*\-]+$`)
var ValidEmail = regexp.MustCompile(`^([0-9a-zA-Z]([-.\w]*[0-9a-zA-Z])*@(([0-9a-zA-Z])+([-\w]*[0-9a-zA-Z])*\.)+[a-zA-Z]{2,9})$`)
var ValidTel = regexp.MustCompile(`^([\(]{1}[0-9]{2,3}[\)]{1}[\.| |\-]{0,1}|^[0-9]{3,4}[\.|\-| ]?)?[0-9]{3,4}(\.|\-| )?[0-9]{3,4}$`)
var ValidRfc = regexp.MustCompile(`^([A-Z&Ññ]{3}|[A-Z][AEIOU][A-Z]{2})\d{2}((01|03|05|07|08|10|12)(0[1-9]|[12]\d|3[01])|02(0[1-9]|[12]\d)|(04|06|09|11)(0[1-9]|[12]\d|30))([A-Z0-9]{2}[0-9A])?$`)
var ValidUrl = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \?=.-]*)*\/?$`)
var ValidPercent = regexp.MustCompile(`^-?[0-9]{0,2}(\.[0-9]{1,2})?$|^-?(100)(\.[0]{1,2})?$`)
var ValidPrice = regexp.MustCompile(`^(\d{1,3},(\d{3}')*\d{3}(\.\d{1,3})?|\d{1,3}(\.\d{2})?)$`)
var ValidSearchData = regexp.MustCompile(`^[a-zA-ZáéíóúAÉÍÓÚÑñäëïöü]+$`)
var ValidID = regexp.MustCompile(`^[a-zA-Z]+$`)
var ValidAlfa = regexp.MustCompile(`^[a-z]+$`)

