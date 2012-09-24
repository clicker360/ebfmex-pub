package site

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
	"strconv"
    "net/http"
    "model"
    "fmt"
)


func init() {
    http.HandleFunc("/imgxp", imgxport)

}

func imgxport(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    if u := user.Current(c); u == nil {
		return
	}

	n, _ := strconv.Atoi(r.FormValue("n"))
	//o, _ := strconv.Atoi(r.FormValue("o"))
    //q := datastore.NewQuery("Cta").Order("FechaHora").Limit(n).Offset(o)
    q := datastore.NewQuery("Cta").Order("FechaHora").Limit(n)
    regdata := make([]model.Cta,0,n)

    if _, err := q.GetAll(c, &regdata); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\">\n");
	fmt.Fprintf(w, "<html>\n");
	fmt.Fprintf(w, "<head>\n");
	fmt.Fprintf(w, "<META http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n");
	fmt.Fprintf(w, "<title>Untitled Document</title>\n");
	fmt.Fprintf(w, "<style type=\"text/css\">\n");
	fmt.Fprintf(w, ".divCenter{ width:95%; margin-left:auto; margin-right:auto; text-align:center;}\n");
	fmt.Fprintf(w, ".imgCont{ width:60px; height:60px; float:left; margin:10px;}\n");
	fmt.Fprintf(w, ".imgCont img { width:100%;}\n");
	fmt.Fprintf(w, "</style>\n");
	fmt.Fprintf(w, "</head>\n");
	fmt.Fprintf(w, "<body>\n");
	fmt.Fprintf(w, "<div class=\"divCenter\">\n");

	for _, cta := range regdata {
		q2 := datastore.NewQuery("Empresa").Ancestor(cta.Key(c))
		for cursor := q2.Run(c); ; {
			var emp model.Empresa
			_, err := cursor.Next(&emp)
			if err == datastore.Done  {
				break
			}
			imgq := datastore.NewQuery("EmpLogo").Filter("IdEmp =", emp.IdEmp)
			for imgcur := imgq.Run(c); ; {
				var img model.Image
				_, err := imgcur.Next(&img)
				if err == datastore.Done  {
					break
				}
				if(img.Data != nil) {
					fmt.Fprintf(w, "<div class=\"imgCont\"><img src=\"/simg?id=%s\" /></div>\n", emp.IdEmp)
				}
			}
		}

	}
	fmt.Fprintf(w, "</div>\n</body>\n");
	fmt.Fprintf(w, "</html>");
}

