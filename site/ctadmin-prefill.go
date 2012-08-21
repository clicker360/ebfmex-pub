package site

import (
    "appengine"
    "net/http"
    "sess"
)

func init() {
    http.HandleFunc("/registro-prefill", registroprefill)
}

func registroprefill(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        if _, ok := sess.IsSess(w, r, c); !ok {
                var fd FormDataCta
		fd.Nombre = r.FormValue("nombre")
		fd.Apellidos = r.FormValue("apellidos")
		fd.Email = r.FormValue("email")
                tc := make(map[string]interface{})
                //tc["Sess"] = s
                tc["FormDataCta"] = fd
                registroTpl.ExecuteTemplate(w, "cta", tc)
                return
        } else {
                http.Redirect(w, r, "/dash", http.StatusFound)
        }
}
