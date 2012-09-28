package site

import (
    //"appengine"
	//"html/template"
	"encoding/json"
    "net/http"
	"fmt"
	//"model"
)

func init() {
    http.HandleFunc("/ws_sucursales", wsSucursales)
    http.HandleFunc("/wsofertas", TopOfertas)
    http.HandleFunc("/wsoferta", Oferta)
    http.HandleFunc("/wsofertaxc", OfertasPorCategoria)
    http.HandleFunc("/wsofertaxp", OfertasPorPalabraClave)
    http.HandleFunc("/wsfaq", PreguntasFrecuentes)
    //http.HandleFunc("/consejos", Consejos)
    //http.HandleFunc("/capsulas", Capsulas)
}

// Envía el ID de la oferta que se consulta
// Trae un Objeto
type Coord struct{
	Params CoordData `json:"Params"`
}
type CoordData struct{
	Latitud		string `json:"Latitud"`
	Longitud	string `json:"Longitud"`
	Distancia	string `json:"Distancia"`
}

func wsSucursales(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	data := new(Coord)
	if err := decoder.Decode(&data); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//c.Errorf("Error decoding data: %s", err)
	}
	//lat := data.Params.Latitud;
	//lng := data.Params.Longitud;
	//rad := data.Params.Distancia

	//w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{ lat:\"12\", lng:\"13\", rad:\"14\" }")
//	fmt.Fprint(w, "{ lat:\"%s\", lng:\"%s\", rad:\"%s\" }", lat, lng, rad)

/*
	if empresa := model.GetEmpresa(c, r.FormValue("IdEmp")); empresa != nil {
		= listSuc(c, empresa.IdEmp)
	}
	*/
}

func Oferta(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Accept-Charset","utf-8");
	w.Header().Set("Content-Type", "application/json")
    var json =`{"id":1, "empresa":{"id": 1, "nombre":"Marti" }, "tipo_oferta":1, "oferta":"tennis Nike", "descuento":25, "descripcion":"Todos los tenis NIKE modelo X, 25% de descuento", "distancia":384, "sucursales": [ { "id": 1, "lat": "19.123456", "lon": "-99.123456" }, { "id": 2, "lat": "19.123456", "lon": "-99.123456" }, { "id": 3, "lat": "19.123456", "lon": "-99.123456" } ], "ofertas_relacionadas": [ { "id": 33, "oferta": "tennis 2 x 1" }, { "id": 34, "oferta": "tennis 2 x 1" } ], "url":"http://www.elbuenfin.org", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv" }`
	fmt.Fprintf(w, "%s", json)
}

// Recibe latitud, longitud y distancia, ejemplo : PARAMS {“latitud”:”3123”,”longitud”:”1231”,”distancia”:”13123”}
//
// Devuelve arreglo de 10 últimas ofertas dependiento de las coordenadas
// Las ofertas son de preferencia las últimas agregadas a las sucursales cercanas
// En caso de no haber sucursales cercanas en esa ubicación se debe de aumentar el radio hasta competar 10.
// En caso de no obtener nada devolver código de error -1
// Las ofertas relacionadas se limitan a 3. (Ofertas relacionadas son 3 ofertas de la misma empresa y devolver
// en orden de cercania
// En caso de no enviar PARAMS regresar 10 ofertas random
// Tipo de Oferta: 1 si es descuento, 2 si es promoción. Si es descuento se debe incluir descuento,
//                 si es promoción se debe poner "promoción" y su tipo "2x1" "meses sin intereses", etc.

func TopOfertas(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Accept-Charset","utf-8");
	w.Header().Set("Content-Type", "application/json")
	var json = `{ "ofertas": [ { "id": 1, "empresa": { "id": 1, "nombre": "Marti" }, "tipo_oferta":1, "oferta": "tennis Nike", "descuento":25, "descripcion": "Todos los tenis Nike modelo X, 25% de descuento", "distancia":384, "sucursales": [ { "id": 1, "lat": "19.123456", "lon": "-99.123456" }, { "id": 2, "lat": "19.123456", "lon": "-99.123456" }, { "id": 3, "lat": "19.123456", "lon": "-99.123456" } ], "ofertas_relacionadas": [ { "id": 35, "oferta": "tennis para fútbol 0 intereses" }, { "id": 34, "oferta": "tennis 2 x 1" } ], "url":"http://www.elbuenfin.org", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv" }, { "id": 2, "empresa": { "id": 2, "nombre": "Emico" }, "tipo_oferta":1, "oferta": "Zapato de vestir, rebajado", "descuento":25, "descripcion": "Toda la linea para caballero 30% de descuento en efectivo", "distancia":38, "sucursales": [ { "id": 1, "lat": "19.408681", "lon": "-99.144552" }, { "id": 2, "lat": "19.401051", "lon": "-99.133072" }, { "id": 3, "lat": "19.37411", "lon": "-99.154251" } ], "ofertas_relacionadas": [ { "id": 33, "oferta": "tennis 2 x 1" }, { "id": 34, "oferta": "tennis 2 x 1" } ], "url":"http://www.elbuenfin.org", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv" }, { "id": 1, "empresa": { "id": 3, "nombre": "Reebok" }, "tipo_oferta":2, "oferta": "tennis para fútbol 0 intereses", "descripcion": "Todos los tenis para fútbol a 6 meses sin intereses", "promocion":"6 meses sin intereses": "distancia":384, "sucursales": [ { "id": 1, "lat": "19.342144", "lon": "-99.144959" }, { "id": 2, "lat": "19.391922", "lon": "-99.19561" } ], "ofertas_relacionadas": [ { "id": 33, "oferta": "tennis 2 x 1" }, { "id": 34, "oferta": "tennis 2 x 1" } ], "url":"http://www.elbuenfin.org", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv" } ] }`
	fmt.Fprintf(w, "%s", json)
}

// Recibe categoría a buscar, status 0 si es presencial y 1 si es en línea, start (entero para paginación),
// Latitud, longitud o entida, si la entidad es 0 es sobre todas las entidades.
//
// Regresa un arreglo de objetos oferta ordenadas por distancia si es que se mando la latitud y longitud,
// De caso contrario se ordenarían alfabéticamente.
// En caso de que se envíe lat y lng se debe devolver el campo distancia
func OfertasPorCategoria(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Accept-Charset","utf-8");
	w.Header().Set("Content-Type", "application/json")
	var json = `{ "ofertas": [ { "id": 33, "oferta": "tennis 2 x 1", "distancia":"1", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv", "empresa": { "id": 2, "nombre": "Marti" } }, { "id": 34, "oferta": "Zapato 2 x 1 1/2" , "distancia":"1.2", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv", "empresa": { "id": 1, "nombre": "Emico" } } ] }`
	fmt.Fprintf(w, "%s", json)
}

// Se envia keyword, start (que es un entero para indicar el indice de paginacion), latitud y longitud o estado.
// Regresa un arreglo de objetos oferta ordenadas por distancia si es que se mando la latitud y longitud, 
// de caso contrario se ordenarian alfabéticamente
func OfertasPorPalabraClave(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Accept-Charset","utf-8");
	w.Header().Set("Content-Type", "application/json")
	var json = `{ "ofertas": [ { "id": 33, "oferta": "tennis 2 x 1", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv", "empresa": { "id": 2, "nombre": "Marti" } }, { "id": 34, "oferta": "tennis 2x1 1/2", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv", "empresa": { "id": 1, "nombre": "Emico" } } ] }`
	fmt.Fprintf(w, "%s", json)
}

func PreguntasFrecuentes(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Accept-Charset","utf-8");
	w.Header().Set("Content-Type", "application/json")
	var json = `{ "preguntas": [ { "pregunta": "¿Que es el buen fin?", "respuesta": "un programa de ofertas" }, { "pregunta": "¿Que es el buen fin?", "respuesta": "un programa de ofertas" } ] }`
	fmt.Fprintf(w, "%s", json)
}

/*
func TopOfertas(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
	_, ok := sess.IsSess(w, r, c)
	if ok {
		if entidad, err := model.GetEntidad(c, r.FormValue("CveEnt")); err == nil {
			if municipios, _ := entidad.GetMunicipios(c); err == nil {
				//Despliega municipios
				tpl, _ := template.New("Mun").Parse(OptionTpl)
				fmt.Fprintf(w, `<select name="DirMun" class="last" id="MunSelector" onchange="locateAddress();">`)
				for _, m := range *municipios {
					//fmt.Fprintf(w, "mun: %s, %s", r.FormValue("DirEnt"), m)
					// Ojo: ver porqué los repite con datos en blanco
					// El if es para corregir la bronca temporalmente
					if(m.Municipio != "") {
						if (m.CveMun == r.FormValue("CveMun")) {
							m.Selected = "selected"
						}
						tpl.Execute(w, m)
					}
				}
				fmt.Fprintf(w, `</select>`)
			}
		}
	}
	return
}
*/

//const OptionTpl = `<option value="{{.CveMun}}" {{if .Selected}}selected="{{.Selected}}"{{end}}>{{.Municipio}}</option>`
