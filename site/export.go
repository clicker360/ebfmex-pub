package site

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
    "net/http"
    "model"
    "fmt"
)

func init() {
    http.HandleFunc("/registros.csv", registroCsv)
}

func registroCsv(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    if u := user.Current(c); u == nil {
		return
	}
    q := datastore.NewQuery("Cta").Order("FechaHora")
    regdata := make([]model.Cta,0,500)

    if _, err := q.GetAll(c, &regdata); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-type", "application/octet-stream");
	w.Header().Set("Content-Disposition", "attachment; filename=\"reportecta.csv\"");
	w.Header().Set("Accept-Charset","utf-8");

	fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
	"cta.Nombre", "cta.Apellidos", "cta.Puesto", "cta.Email", "cta.EmailAlt", "cta.Pass", "cta.Tel", "cta.Cel", "cta.FechaHora", "cta.CodigoCfm", "cta.Status",
	"IdEmp", "RFC", "Nombre Empresa", "Razon Social", "Dir.Calle", "Dir.Colonia", "Dir.Entidad", "Dir.Municipio", "Dir.Cp", "Dir.Número Suc",
	"Organiso Emp", "Otro Organismo", "Reg Org. Empresarial", "Url", "PartLinea", "ExpComer", "Descripción", "FechaHora Alta Emp.","emp.Status")
	for _, cta := range regdata {
		q2 := datastore.NewQuery("Empresa").Ancestor(cta.Key(c))
		for cursor := q2.Run(c); ; {
			var emp model.Empresa
			_, err := cursor.Next(&emp)
			if err == datastore.Done  {
				break
			}
			fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%t\n",
			cta.Nombre, cta.Apellidos, cta.Puesto, cta.Email, cta.EmailAlt, cta.Pass, cta.Tel, cta.Cel, cta.FechaHora, cta.CodigoCfm, cta.Status,
			emp.IdEmp, emp.RFC, emp.Nombre, emp.RazonSoc, emp.DirCalle, emp.DirCol, emp.DirEnt, emp.DirMun, emp.DirCp, emp.NumSuc,
			emp.OrgEmp, emp.OrgEmpOtro, emp.OrgEmpReg, emp.Url, emp.PartLinea, emp.ExpComer, emp.Desc, emp.FechaHora, emp.Status)
		}

	}
}

