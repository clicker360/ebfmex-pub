package site

import (
    	"appengine"
    	"appengine/datastore"
//    "appengine/user"
    	"net/http"
        "fmt"
        "model"
	"strconv"
	//"strings"
)

func init() {
    http.HandleFunc("/api/cerca", Cerca)
}

func Cerca(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Sucursal")
	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	long, _ := strconv.ParseFloat(r.FormValue("long"), 64)
	rad, _ := strconv.ParseFloat((r.FormValue("rad")), 64)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Cerca de %f, %f - Radius: %f\n", lat, long, rad);
	for i := q.Run(c); ; {
		var s model.Sucursal
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
}
