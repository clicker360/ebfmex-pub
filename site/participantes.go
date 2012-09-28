package site

import (
    "appengine"
    "appengine/datastore"
    "appengine/user"
	//"strconv"
    "net/http"
	//"sort"
    "model"
	"sortutil"
    "fmt"
)


func init() {
    http.HandleFunc("/participantes", ShowEmpresas)

}

func ShowEmpresas(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    if u := user.Current(c); u == nil {
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


	/*
		Loop para recorrer todas las empresas 
	*/
	prefix := r.FormValue("prefix")
	date := r.FormValue("date")
	//o, _ := strconv.Atoi(r.FormValue("o"))
    //q := datastore.NewQuery("Cta").Order("FechaHora").Limit(n).Offset(o)
    //q := datastore.NewQuery("Empresa").Order("FechaHora").Limit(n)
	q := datastore.NewQuery("Empresa").Filter("Nombre >=", prefix).Filter("Nombre <", prefix+"\ufffd")
	em, _ := q.Count(c)
	empresas := make([]model.Empresa, 0, em)
	if _, err := q.GetAll(c, &empresas); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return
		}
	}

	sortutil.AscByField(empresas, "Nombre")
	if date != "" {
		sortutil.AscByField(empresas, "FechaHora")
	}
	for k, _ := range empresas {
		imgq := datastore.NewQuery("EmpLogo").Filter("IdEmp =", empresas[k].IdEmp)
		for imgcur := imgq.Run(c); ; {
			var img model.Image
			_, err := imgcur.Next(&img)
			if err == datastore.Done  {
				break
			}
			if(img.Data != nil) {
				fmt.Fprintf(w, "<div class=\"imgCont\"><img src=\"/simg?id=%s\" />%s</div>\n", empresas[k].IdEmp, empresas[k].Nombre)
			}
		}

	}
	fmt.Fprintf(w, "</div>\n</body>\n");
	fmt.Fprintf(w, "</html>");
}

