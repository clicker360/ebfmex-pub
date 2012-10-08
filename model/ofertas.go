package model

import (
    "appengine"
    "appengine/datastore"
	"time"
)

type Oferta struct {
	IdOft       string
	IdEmp       string
	IdCat       int
	Empresa		string
	Oferta		string
	NOferta			string
	Descripcion		string
	NDescripcion	string
	Codigo      string
	Precio      string
	Descuento   string
	Promocion	string
	Enlinea     bool
	Url         string
	Tarjetas    []byte // json
	Meses       string
	FechaHoraPub    time.Time
	StatusPub   bool
	FechaHora   time.Time
	Image	[]byte
	ImageA	[]byte
	ImageB	[]byte
	Sizepx	int
	Sizepy	int
	SizeApx	int
	SizeApy	int
	SizeBpx	int
	SizeBpy	int
}

type OfertaSucursal struct {
	IdOft       string
	IdEmp       string
	IdSuc       string
	Sucursal    string
	Lat         float64
	Lng         float64
	Empresa     string
	Oferta      string
	NOferta		string
	Descripcion		string
	NDescripcion	string
	Promocion	string
	Precio      string
	Descuento   string
	Url         string
	StatusPub   bool
	FechaHora	time.Time
}

type Categoria struct {
	IdCat       int
	Categoria   string
	Selected	string `datastore:"-"`
}

type OfertaPalabra struct {
	IdOft      string
	IdEmp      string
	Palabra    string
	FechaHora	time.Time
}

type OfertaImage struct {
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
	FechaHora	time.Time
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

func GetOferta(c appengine.Context, id string) *Oferta {
	q := datastore.NewQuery("Oferta").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e Oferta
		_, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		// Regresa la oferta
		return &e
	}
	// Regresa un cascarón
	var e Oferta
	e.IdEmp = "none";
	e.IdOft = "none";
	e.IdCat = 0;
	return &e
}

func PutOferta(c appengine.Context, oferta *Oferta) error {
	_, err := datastore.Put(c, oferta.Key(c), oferta)
	if err != nil {
		return err
	}
	/* 
		relación oferta sucursal 
	*/

	return nil
}

func NewOferta(c appengine.Context, oferta *Oferta) error {
	oferta.IdOft = RandId(14)
    _, err := datastore.Put(c, datastore.NewKey(c, "Oferta", oferta.IdOft, 0, nil), oferta)
	if err != nil {
		return err
	}
	return nil
}

func GetOfertaSucursales(c appengine.Context, idoft string) (*[]OfertaSucursal, error) {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", idoft)
	n, _ := q.Count(c)
	ofersuc := make([]OfertaSucursal, 0, n)
	if _, err := q.GetAll(c, &ofersuc); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return &ofersuc, nil
}

func GetOfertaPalabras(c appengine.Context, idoft string, idemp string) *[]OfertaPalabra {
	q := datastore.NewQuery("OfertaPalabra")
	if(idemp != "") {
		q = q.Filter("IdEmp =", idemp)
	} else {
		q = q.Filter("IdOft =", idoft)
	}
	n, _ := q.Count(c)
	op := make([]OfertaPalabra, 0, n)
	if _, err := q.GetAll(c, &op); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil
		}
	}
	return &op
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

func GetCategoria(c appengine.Context, id int) *Categoria {
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
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaSucursal", r.IdOft+ofsuc.IdSuc, 0, r.Key(c)), ofsuc)
	if err != nil {
		return err
	}
	return nil
}

/*
	Llenar primero struct de OfertaPalabra y luego guardar
*/
func (r *Oferta) PutOfertaPalabra(c appengine.Context, op *OfertaPalabra) error {
	_, err := datastore.Put(c, datastore.NewKey(c, "OfertaPalabra", r.IdEmp+op.Palabra, 0, r.Key(c)), op)
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
		if err:= DelOfertaSucursales(c, id); err != nil {
			return err
		}
		if err := DelOfertaPalabras(c, id); err != nil {
			return err
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

/*
	Método para borrar todas las sucursales de una oferta
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
	Método para borrar todas las sucursales de una oferta
*/
func DelOfertaSucursal(c appengine.Context, idoft string, idsuc string) error {
	q := datastore.NewQuery("OfertaSucursal").Filter("IdOft =", idoft).Filter("IdSuc =", idsuc)
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
	Las palabras clave asociadas a una oferta se ponen con idoft="none" todas juntas
*/
func DelOfertaPalabras(c appengine.Context, id string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		e.IdOft = "none";
		/*
			En realidad no se borra ningun entity, solo se desliga la oferta
			La palabra continua perteneciendo a la empresa para uso de las demás
			ofertas
		*/
		_, err = datastore.Put(c, key, &e)
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

/*
	Las palabras clave asociadas a una oferta se borran todas juntas
*/
func RmOfertaPalabras(c appengine.Context, id string) error {
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


/*
	Las palabra clave asociada se borra individualmente
*/
func RmOfertaPalabra(c appengine.Context, id string, palabra string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id).Filter("Palabra =", palabra)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		/*
			Aquí si se borra el entity
		*/
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
	}
	return nil
}

func DelOfertaPalabra(c appengine.Context, id string, palabra string) error {
	q := datastore.NewQuery("OfertaPalabra").Filter("IdOft =", id).Filter("Palabra =", palabra)
	for i := q.Run(c); ; {
		var e OfertaPalabra
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		e.IdOft = "none";
		/*
			En realidad no se borra ningun entity, solo se desliga la oferta
			La palabra continua perteneciendo a la empresa para uso de las demás
			ofertas
		*/
		_, err = datastore.Put(c, key, &e)
	}
	return nil
}



