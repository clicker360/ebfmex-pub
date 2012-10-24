package backend

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"strconv"
	"fmt"
)

func init() {
    http.HandleFunc("/backend/updatesearch", fetchUpdateSearch)
}

func fetchUpdateSearch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	/* 
	 * En lo que encontramos una manera digna de ejecutar el cron, se
	 * podría meter una llave en el memcaché y crearla en el administrador
	 * O bien en el datastore.
	 */
	var m int
	m, _ = strconv.Atoi(r.FormValue("m"))
	if m < 30 {
		m = 30
	}
	if r.FormValue("c") == "ZWJmbWV4LXB1YnIeCxISX0FoQWRtaW5Yc3JmVG9rZW5fIgZfWFNSRl8M" {
		descurl := fmt.Sprintf( "http://movil.%s.appspot.com/backend/updatesearch?minutes=%d", appengine.AppID(c), m)
		//descurl := fmt.Sprintf("http://clicker360.com")
		_, err := client.Get(descurl)
		if err != nil {
			c.Errorf("updatesearch urlfetch %v client: %v", appengine.AppID(c), err)
		}
		c.Infof("updatesearch urlfetch client %v minutes=%d", appengine.AppID(c), m)
	}
	return
}
