package model

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
	IdEmp		string
	Nombre		string
	Tel			string
	DirCalle	string
	DirCol		string
	DirEnt		string
	DirMun		string
	DirCp		string
	GeoUrl		string
	Geo1		string
	Geo2		string
	Geo3		string
	Geo4		string
	FechaHora	time.Time
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

type Image struct {
	Data	[]byte
	IdEmp	string
	IdImg	string
	Kind	string
	Name	string
	Desc	string
	Sizepx	int
	Sizepy	int
	Url		string
	Type	string
	Sp1		string
	Sp2		string
	Sp3		string
	Sp4		string
	Np1		int
	Np2		int
	Np3		int
	Np4		int
}

// Interfaces
//type Lister interface {
//	ListAll() 
//}

func (r *Cta) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Cta", r.Email, 0, nil)
}

func GetCta(c appengine.Context, email string) (*Cta, error) {
	ua := &Cta{ Email: email, }
	err := datastore.Get(c, ua.Key(c), ua)
	if err == datastore.ErrNoSuchEntity {
		return ua, err
		//_, err = datastore.Put(c, ua.Key(), ua)
	}
	return ua, nil
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

func (r *Cta) GetEmpresa(c appengine.Context, id string) (*Empresa, error) {
	e := &Empresa{ IdEmp: id }
	err := datastore.Get(c, datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c)), e)
	if err == datastore.ErrNoSuchEntity {
		return e, err
	}
	return e, err
}

func (r *Cta) PutEmpresa(c appengine.Context, e *Empresa) (*Empresa, error) {
	_, err := datastore.Put(c, datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c)), e)
	if err != nil {
		return nil, err
	}
	return e, err
}

func (r *Cta) NewEmpresa(c appengine.Context, e *Empresa) (*Empresa, error) {
	e.IdEmp = RandId(12)
    _, err := datastore.Put(c, datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c)), e)
	if err != nil {
		return nil, err
	}
	return e, err
}

func (r *Cta) DelEmpresa(c appengine.Context, id string) error {
	_ = DelImg(c, "EmpLogo", id)
    if err := datastore.Delete(c, datastore.NewKey(c, "Empresa", id, 0, r.Key(c))); err != nil {
		return err
	}
	return nil
}

// Métodos de Empresa
func GetEmpresa(c appengine.Context, id string) (*Empresa) {
	q := datastore.NewQuery("Empresa").Filter("IdEmp =", id)
	for i := q.Run(c); ; {
		var e Empresa
		_, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		return &e
	}
	return nil
}

func (e *Empresa) Key(c appengine.Context) *datastore.Key {
	//return datastore.NewKey(c, "Empresa", e.IdEmp, 0, nil)
	q := datastore.NewQuery("Empresa").Filter("IdEmp =", e.IdEmp)
	for i := q.Run(c); ; {
		var e Empresa
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		return key
	}
	return nil
}

func (e *Empresa) PutSuc(c appengine.Context, s *Sucursal) (*Sucursal, error) {
	if(s.IdSuc == "none" || s.IdSuc == "") {
		s.IdSuc = RandId(14)
	}
	parentKey := e.Key(c)
    _, err := datastore.Put(c, datastore.NewKey(c, "Sucursal", s.IdSuc, 0, parentKey), s)
	if err != nil {
		return nil, err
	}
	return s, err
}

// Métodos de Sucursal
func GetSuc(c appengine.Context, id string) (*Sucursal) {
	q := datastore.NewQuery("Sucursal").Filter("IdSuc =", id)
	for i := q.Run(c); ; {
		var e Sucursal
		_, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		// Regresa la sucursal
		return &e
	}
	// Regresa un cascarón
	var e Sucursal
	e.IdSuc = "none";
	e.IdEmp = "none";
	return &e
}

func DelSuc(c appengine.Context, id string) error {
	q := datastore.NewQuery("Sucursal").Filter("IdSuc =", id)
	for i := q.Run(c); ; {
		var e Sucursal
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

// Métodos de Entidad
func (e *Entidad) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Entidad", e.CveEnt, 0, nil)
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

func GetEntidad(c appengine.Context, cveent string) (*Entidad, error) {
	e := &Entidad{ CveEnt: cveent }
	err := datastore.Get(c, e.Key(c), e)
	if err == datastore.ErrNoSuchEntity {
		return nil, err
	}
	return e, nil
}



// Métodos de Municipio
func (m *Municipio) Parent(c appengine.Context) *Entidad {
	e, _ := GetEntidad(c, m.CveEnt)
	return e
}

func (m *Municipio) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Municipio", m.CveMun, 0, m.Parent(c).Key(c))
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

// Métodos de Imagen
// Obtiene la llave de una imagen
func (i *Image) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, i.Kind, i.IdImg, 0, nil)
}

// Borra una imagen
func DelImg(c appengine.Context, kind string, id string) error {
    if err := datastore.Delete(c, datastore.NewKey(c, kind, id, 0, nil)); err != nil {
		return err
	}
	return nil
}

// Guarda Imagen modificada
func PutLogo(c appengine.Context, i *Image) (*datastore.Key, error) {
	key, err := datastore.Put(c, datastore.NewKey(c, i.Kind, i.IdEmp, 0, nil), i)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GetLogo(c appengine.Context, idemp string) (*Image, error) {
	i := &Image{ IdEmp: idemp, Kind: "EmpLogo" }
	// Para el logo sólo se utiliza la llave IdEmp
	err := datastore.Get(c, datastore.NewKey(c, i.Kind, i.IdEmp, 0, nil), i)
	if err == datastore.ErrNoSuchEntity {
		// Se crea el spaceholder del logo
		_, err := PutLogo(c, i)
		if err != nil {
			return nil, err
		}
	}
	return i, err
}

// Obtiene una imagen
func GetImg(c appengine.Context, id string) (*Image, error) {
	i := &Image{ IdImg: id }
	err := datastore.Get(c, i.Key(c), i)
	if err == datastore.ErrNoSuchEntity {
		//_, err = datastore.Put(c, ua.Key(), ua)
		return i, err
	}
	return i, err
}


