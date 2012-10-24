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
	var minutes int
	minutes, _ = strconv.Atoi(r.FormValue("minutes"))
	if minutes < 30 {
		minutes = 30
	}
	descurl := fmt.Sprintf( "http://movil.%s.appspot.com/backend/updatesearch?minutes=%d", appengine.AppID(c), minutes)
	_, err := client.Get(descurl)
	if err != nil {
		c.Errorf("updatesearch urlfetch client: %v", err)
	}
	c.Errorf("updatesearch urlfetch client minutes=%d", minutes)
	return nil
}
