/*
	GET         /ent
	GET         /ent/12
	GET         /munsde/12


	RESTful API resources

	GET         /api.json
	GET         /api/01.json

	RespondWithNotFound() looks like:

	/api/no-ruch-resource.json


	See JSONP in action with these URLs:

	http://localhost:8080/api.json?callback=MyFunc
	http://localhost:8080/api.json?callback=MyFunc&context=123
*/
package site

import (
    //"appengine"
	//"net/http"
	"goweb"
	"fmt"
)

/*
	API Ofertas
*/
type OfertaCtl struct {}
func (mc *OfertaCtl) HandleRequest(cx *goweb.Context) {
	/*
	ac := appengine.NewContext(cx.Request)
	if oferta, err := GetOferta(ac, cx.PathParams["id"]); err != nil {
		cx.RespondWithNotFound()
	} else {
		cx.RespondWithData(&oferta)
	}
	*/
    var json =`{"id":1, "empresa":{"id": 1,"nombre":"Marti"}, "tipo_oferta":1 "oferta":"tennis 2 x 1", "descripcion":"Descuento en NIKE", "distancia":384, "sucursales": [ { "id": 1, "lat": "19.123456", "long": "-99.123456" }, { "id": 2, "lat": "19.123456", "long": "-99.123456" }, { "id": 3, "lat": "19.123456", "long": "-99.123456" } ], "ofertas_relacionadas": [ { "id": 33, "oferta": "tennis 2 x 1" }, { "id": 34, "oferta": "tennis 2 x 1" } ], "url":"http://www.elbuenfin.org", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv" }`
	w := cx.ResponseWriter
	fmt.Fprintf(w, "%s", json)
}

type TopOfertasCtl struct {}
func (mc *TopOfertasCtl) HandleRequest(cx *goweb.Context) {
	/*
	ac := appengine.NewContext(cx.Request)
	if ofertas, err := GetTopOfertas(ac, cx.PathParams["lat"], cx.PathParams["lng"], cx.PathParams["rad"]); err != nil {
		cx.RespondWithNotFound()
	} else {
		cx.RespondWithData(&ofertas)
	}
	*/
	var json = `{ "ofertas": [ { "id": 1, "empresa": { "id": 1, "nombre": "Marti" }, "tipo_oferta":1 "oferta": "tennis 2 x 1", "descripcion": "Descuento en NIKE", "distancia":384, "sucursales": [ { "id": 1, "lat": "19.123456", "long": "-99.123456" }, { "id": 2, "lat": "19.123456", "long": "-99.123456" }, { "id": 3, "lat": "19.123456", "long": "-99.123456" } ], "ofertas_relacionadas": [ { "id": 33, "oferta": "tennis 2 x 1" }, { "id": 34, "oferta": "tennis 2 x 1" } ], "url":"http://www.elbuenfin.org", "url_image":"http://www.elbuenfin.org/simg?id=tmshqympnnvv" } ] }`
	w := cx.ResponseWriter
	fmt.Fprintf(w, "%s", json)
}

type OfertasPorCategoriaCtl struct {}
//func /mc *OfertaPorCategoriaCtl
/*
	
	The main function will register the relevant controllers
	and start the web server

*/
func init() {

	/*
		API controller
	*/
	var ofertaCtl *OfertaCtl = new(OfertaCtl)
	goweb.Map("/wsroferta/{id}", ofertaCtl)

	var topOfertasCtl *TopOfertasCtl = new(TopOfertasCtl)
	goweb.Map("/wsrofertas/{lat}/{lng}/{rad}", topOfertasCtl)

}
