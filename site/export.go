package site

import (
    "appengine"
    "appengine/datastore"
    "net/http"
    "sess"
    "html/template"
    "time"
    "model"
    //"fmt"
)

type Cta struct {
        Folio                   int32
        Nombre                  string
        Apellidos               string
        Puesto                  string
        Email                   string
        EmailAlt                string
        Pass                    string
        Tel                             string
        Cel                             string
        FechaHora               time.Time
        UsuarioInt              string
        CodigoCfm               string
        Status                  bool
}

type Empresa struct {
        IdEmp           string
        Folio           int32
        RFC                     string
        Nombre          string
        RazonSoc        string
        DirCalle        string
        DirCol          string
        DirEnt          string
        DirMun          string
        DirCp           string
        NumSuc          string
        OrgEmp          string
        OrgEmpOtro      string
        OrgEmpReg       string
        //Entidades     []Entidad
        Url                     string
        Benef           int
        PartLinea       int
        ExpComer        int
        Desc            string
        FechaHora       time.Time
        Status          bool
}

type CtaEmpresa struct {
        Folio                   int32
        Nombre                  string
        Apellidos               string
        Puesto                  string
        Email                   string
        EmailAlt                string
        Pass                    string
        Tel                             string
        Cel                             string
        FechaHora               time.Time
        UsuarioInt              string
        CodigoCfm               string
        Status                  bool
        IdEmp           string
        FolioE           int32
        RFC                     string
        NombreE          string
        RazonSoc        string
        DirCalle        string
        DirCol          string
        DirEnt          string
        DirMun          string
        DirCp           string
        NumSuc          string
        OrgEmp          string
        OrgEmpOtro      string
        OrgEmpReg       string
        //Entidades     []Entidad
        Url                     string
        Benef           int
        PartLinea       int
        ExpComer        int
        Desc            string
        FechaHoraE       time.Time
        StatusE          bool
}

func init() {
    http.HandleFunc("/registro-export", registroExport)
    http.HandleFunc("/registros.csv", registroCsv)
}

func registroExport(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        if _, ok := sess.IsSess(w, r, c); !ok {
                tc := make(map[string]interface{})
                exportTpl.ExecuteTemplate(w, "cta", tc)
                return
        } else {
		return
        }
}

func registroCsv(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Cta").Order("FechaHora")
    regdata := make([]model.Cta,0,100)

    if _, err := q.GetAll(c, &regdata); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


	w.Header().Set("Content-Type", "text/plain")

	for _, value := range regdata {
		q2 := datastore.NewQuery("Empresa").Ancestor(value.Key(c)).Limit(50)
		for cursor := q2.Run(c); ; {
			var e Empresa
			_, err := cursor.Next(&e)
			if err == datastore.Done  {
				break
			}
			// llenar linea de csv
			ce := CtaEmpresa {
        			Folio: value.Folio,
			        Nombre: value.Nombre,
			        Apellidos: value.Apellidos,
			        Puesto: value.Puesto,
			        Email: value.Email,
			        EmailAlt: value.EmailAlt,
			        Pass: value.Pass,
			        Tel: value.Tel,
			        Cel: value.Cel,
			        FechaHora: value.FechaHora,
			        UsuarioInt: value.UsuarioInt,
			        CodigoCfm: value.CodigoCfm,
			        Status: value.Status,
			        IdEmp: e.IdEmp,
			        FolioE: e.Folio,
			        RFC: e.RFC,
			        NombreE: e.Nombre,
			        RazonSoc: e.RazonSoc,
			        DirCalle: e.DirCalle,
			        DirCol: e.DirCol,
			        DirEnt: e.DirEnt,
			        DirMun: e.DirMun,
			        DirCp: e.DirCp,
			        NumSuc: e.NumSuc,
			        OrgEmp: e.OrgEmp,
			        OrgEmpOtro: e.OrgEmpOtro,
			        OrgEmpReg: e.OrgEmpReg,
			        Url: e.Url,
			        Benef: e.Benef,
			        PartLinea: e.PartLinea,
			        ExpComer: e.ExpComer,
			        Desc: e.Desc,
			        FechaHoraE: e.FechaHora,
			        StatusE: e.Status,
			}
			//fmt.Fprintf(w, "%s", value.Nombre)
			cuentasCsvTpl.Execute(w, ce)
		}

	}

    //if err := registrosCsvTpl.Execute(w, regdata); err != nil {
//	return
    //}
}

var exportTpl = template.Must(template.ParseFiles("templates/export.html"))
var registrosCsvTpl = template.Must(template.ParseFiles("templates/registros.csv"))
var cuentasCsvTpl = template.Must(template.ParseFiles("templates/cuentas.csv"))
