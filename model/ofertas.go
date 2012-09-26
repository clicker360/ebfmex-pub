package model

import (
    "appengine"
    "appengine/datastore"
	"time"
)

type Oferta struct {
	IdOft       string
	IdEmp       string
	IdCat       string
	Empresa		string
	Oferta      string
	Descripcion	string
	Codigo      string
	Precio      string
	Descuento   string
	Enlinea     bool
	Url         string
	Tarjetas    string // Texto separado por comas
	Meses       string
	FechaHoraPub    time.Time
	StatusPub   bool
	FechaHora   time.Time
}

type OfertaSucursal struct {
	IdOft       string
	IdEmp       string
	IdSuc       string
	Sucursal    string
	lat         float64
	lng         float64
	Empresa     string
	Oferta      string
	Precio      string
	Descuento   string
	Url         string
	StatusPub   bool
}

type Categoria struct {
	IdCat       int
	Categoria   string
}

type OfertaPalabra struct {
	IdSuc      string
	IdOft      string
	Palabra    string
}

func (r *Oferta) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Oferta", r.IdOft, 0, nil)
}

func (r *Oferta) DelOferta(c appengine.Context) error {
    if err := datastore.Delete(c, r.Key(c)); err != nil {
		return err
	}
	return nil
}

func GetOferta(c appengine.Context, id string) (*Oferta, error) {
	oferta := &Oferta{ IdOft: id, }
	err := datastore.Get(c, oferta.Key(c), oferta)
	if err == datastore.ErrNoSuchEntity {
		return nil, err
	}
	return oferta, nil
}

func PutOferta(c appengine.Context, oferta *Oferta) (*Oferta, error) {
	_, err := datastore.Put(c, oferta.Key(c), oferta)
	if err != nil {
		return nil, err
	}
	return oferta, err
}

func GetOfertaSucursales(c appengine.Context, idoft string) (*[]OfertaSucursal, error) {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", idoft)
	n, _ := q.Count(c)
	ofersuc := make([]OfertaSucursal, n)
	if _, err := q.GetAll(c, &ofersuc); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return &ofersuc, nil
}

func GetOfertaSucursalesGeo(c appengine.Context, lat string, lng string, rad string) (*Sucursal, error) {
	/*
	q := datastore.NewQuery("Sucursal")
	for i := q.Run(c); ; {
		var s Sucursal
        _, err := i.Next(&s)
		if err == datastore.Done {
			break
        }
		geo1, _ := strconv.ParseFloat(s.Geo1, 64)
		geo2, _ := strconv.ParseFloat(s.Geo2, 64)
		sqdist := (lat - geo1) * (lat - geo1)  + (long - geo2) * (long - geo2);
		if ( sqdist <= rad * rad) {
			fmt.Fprintf(w, "lat, long: %s, %s\n", s.Geo1, s.Geo2);
		}
	}
	*/
	return nil,nil
}

func GetCategoria(c appengine.Context, id string) *Categoria {
	q := datastore.NewQuery("Categoria").Filter("IdCat =", id).Limit(1)
	for i := q.Run(c); ; {
		var c Categoria
		_, err := i.Next(&c)
		if err == datastore.Done {
			break
		}
		return &c
	}
	return nil
}

/*
	Llenar primero struct de OfertaSucursal y luego guardar
*/
func (r *Oferta) PutOfertaSucursal(c appengine.Context, ofsuc *OfertaSucursal) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaSucursal", "", 0, r.Key(c)), ofsuc)
	if err != nil {
		return err
	}
	return nil
}

/*
	Llenar primero struct de OfertaSucursal y luego guardar
*/
func (r *Oferta) PutOfertaPalabra(c appengine.Context, op *OfertaPalabra) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaPalabra", "", 0, r.Key(c)), op)
	if err != nil {
		return err
	}
	return nil
}

func DelOferta(c appengine.Context, id string) error {
	q := datastore.NewQuery("Oferta").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e Oferta
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

/*
	Las sucursales asociadas a una oferta o se crean todas 
	juntas o se borran todas juntas
*/
func DelOfertaSucursales(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e OfertaSucursal
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

/*
	Las palabras clave asociadas a una oferta o se crean todas 
	juntas o se borran todas juntas
*/
func DelOfertaPalabra(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e OfertaPalabra
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



