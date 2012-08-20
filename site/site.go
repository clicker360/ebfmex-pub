package site

import (
    "appengine"
    "appengine/datastore"
	"net/http"
	"time"
)

func init() {
    http.HandleFunc("/", home)
}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

type Cta struct {
	Folio			int32
	Nombre			string
	Apellidos		string
	Puesto			string
	Email			string
	EmailAlt		string
	Pass			string
	Tel				string
	Cel				string
	FechaHora		time.Time
	UsuarioInt		string
	CodigoCfm		string
	Status			bool
}

type Empresa struct {
	IdEmp		string
	Folio		int32
	RFC			string
	Nombre		string
	RazonSoc	string
	DirCalle	string
	DirCol		string
	DirEnt		string
	DirMun		string
	DirCp		string
	NumSuc		string
	OrgEmp		string
	OrgEmpOtro	string
	OrgEmpReg	string
	//Entidades	[]Entidad
	Url			string
	Benef		int
	PartLinea	int
	ExpComer	int
	Desc		string
	FechaHora	time.Time
	Status		bool
}

type Sucursal struct {
	IdSuc		string
	Nombre		string
	Dir			string
	Tel			string
	GeoUrl		string
	Selected	string
}

type Quest struct {
	PartLinea	int
	ExpComer	int
	Desc		string
}

type Entidad struct {
	CveEnt		string
	Entidad		string
	Abrv		string
	CveCap		string
	Capital		string
	Selected	string
}

type Municipio struct {
	CveEnt		string
	Entidad		string
	Abrv		string
	CveMun		string
	Municipio	string
	CvaCab		string
	Cabecera	string
	Selected	string
}

type Organismo struct {
	Siglas		string
	Nombre		string
	Selected	string
}

// Interfaces
type Lister interface {
	ListAll() 
}

func (r *Cta) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Cta", r.Email, 0, nil)
}

func GetCta(c appengine.Context, email string) (*Cta, error) {
	ua := &Cta{
		Email: email,
	}
	err := datastore.Get(c, ua.Key(c), ua)
	if err == datastore.ErrNoSuchEntity {
		return ua, err
		//_, err = datastore.Put(c, ua.Key(), ua)
	}
	return ua, err
}

func (r *Cta) DelCta(c appengine.Context) error {
    if err := datastore.Delete(c, r.Key(c)); err != nil {
		return err
	}
	return nil
}

func PutCta(c appengine.Context, u *Cta) (*Cta, error) {
	_, err := datastore.Put(c, u.Key(c), u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (r *Cta) KeyEmpresa(c appengine.Context, e *Empresa) *datastore.Key {
	return datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c))
}

func (r *Cta) GetEmpresa(c appengine.Context, id string) (*Empresa, error) {
	e := &Empresa{
		IdEmp: id,
	}
	err := datastore.Get(c, r.KeyEmpresa(c, e), e)
	if err == datastore.ErrNoSuchEntity {
		return e, err
	}
	return e, err
}

func (r *Cta) PutEmpresa(c appengine.Context, e *Empresa) (*Empresa, error) {
	_, err := datastore.Put(c, r.KeyEmpresa(c, e), e)
	if err != nil {
		return nil, err
	}
	return e, err
}

func (r *Cta) AddEmpresa(c appengine.Context, e *Empresa) (*Empresa, error) {
	e.IdEmp = randId(12)
    _, err := datastore.Put(c, r.KeyEmpresa(c, e), e)
	if err != nil {
		return nil, err
	}
	return e, err
}

func (r *Cta) DelEmpresa(c appengine.Context, id string) (*Empresa, error) {
	e := &Empresa{
		IdEmp: id,
	}
    if err := datastore.Delete(c, r.KeyEmpresa(c, e)); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Entidad) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Entidad", e.CveEnt, 0, nil)
}

func (m *Municipio) Parent(c appengine.Context) *Entidad {
	e, _ := GetEntidad(c, m.CveEnt)
	return e
}

func (m *Municipio) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Municipio", m.CveMun, 0, m.Parent(c).Key(c))
}

func GetEntidad(c appengine.Context, cveent string) (*Entidad, error) {
	e := &Entidad{ CveEnt: cveent }
	err := datastore.Get(c, e.Key(c), e)
	if err == datastore.ErrNoSuchEntity {
		return nil, err
	}
	return e, nil
}

func GetMunicipio(c appengine.Context, cvemun string) *Municipio {
	q := datastore.NewQuery("Municipio").Filter("CveMun =", cvemun)
	for i := q.Run(c); ; {
		var m Municipio
		_, err := i.Next(&m)
		if err == datastore.Done {
			break
		}
		return &m
	}
	return nil
}

func (e *Entidad) GetMunicipios(c appengine.Context) (*[]Municipio, error) {
	q := datastore.NewQuery("Municipio").Ancestor(e.Key(c))
	nm, _ := q.Count(c)
	municipios := make([]Municipio, nm)
	if _, err := q.GetAll(c, &municipios); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return &municipios, nil
}


/*
func (r *Empresa) AddSucursal(c appengine.Context, id string, s Sucursal) (*Sucursal, error) {
	key := datastore.NewKey(c, "Sucursal", id, 0, r.Key(c))
	s.IdSuc = id
    //sucursal := &Sucursal{IdSuc: id}
    _, err := datastore.Put(c, key, &s)
	if err != nil {
		return nil, err
	}
	return &s, err
}
*/


/*
                <label>Organismo Empresarial</label>
                <select name="OrgEmp" class="last">
                  <option value="ANTAD">ANTAD</option>
                  <option value="ABM">ABM</option>
                  <option value="AMIPCI">AMIPCI</option>
                  <option value="CCE">CCE</option>
                  <option value="CONCAMIN">CONCAMIN</option>
                  <option value="CONCANACO">CONCANACO</option>
                  <option value="COPARMEX">COPARMEX</option>
                  <option value="NINGUNO">NINGUNO</option>
                <option value="OTRO">OTRO</option>
                </select>

*/
