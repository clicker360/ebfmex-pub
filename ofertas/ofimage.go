// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be fouerrnd in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	"appengine/blobstore"
	"appengine/image"
	"encoding/json"
	"net/http"
	"sess"
	"model"
)

type OfImg struct{
	IdOft		string `json:"idoft"`
	IdBlob		string `json:"idblob"`
	Status		string `json:"errstatus"`
	UploadURL	string `json:"uploadurl"`
}

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {
	http.HandleFunc("/r/ofimgup", handleUpload)
	// ofimg queda fuera del url seguro /r
	//http.HandleFunc("/ofimg", handleServe)
	http.HandleFunc("/ofimg", handleServeImg)

	//http.HandleFunc("/r/ofimgform", handleRoot)
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	blobstore.Send(w, appengine.BlobKey(r.FormValue("id")))
}

func handleServeImg(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") != "none" {
		c := appengine.NewContext(r)
		var imgprops image.ServingURLOptions
		imgprops.Secure = false
		imgprops.Size = 400
		imgprops.Crop = false
		url, _ := image.ServingURL(c, appengine.BlobKey(r.FormValue("id")), &imgprops)
		http.Redirect(w, r, url.String(), http.StatusFound)
	}
	return
}

/* 
 * dejamos esto como referencia
 * El envío de la liga de sesión de upload se genera en ofadm.go
 *
func handleRoot(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
		uploadURL, err := blobstore.UploadURL(c, "/ofimgup", nil)
		if err != nil {
		serveError(c, w, err)
		return
	}
	tc := make(map[string]interface{})
	tc["UploadURL"] = uploadURL
	tc["IdOft"] =  r.FormValue("IdOft")
	w.Header().Set("Content-Type", "text/html")
	rootTemplate.ExecuteTemplate(w, "ofupform", tc)
	return
}
 */

func handleUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var out OfImg
	out.Status = "invalidId"
	out.IdBlob = ""
	if _, ok := sess.IsSess(w, r, c); ok {
		blobs, form, err := blobstore.ParseUpload(r)
		file := blobs["image"]
		out.IdBlob = string(file[0].BlobKey)
		out.IdOft = form.Get("IdOft")
		if err != nil {
			out.Status = "invalidUpload"
			berr := blobstore.Delete(c, file[0].BlobKey)
			model.Check(berr)
		} else {
			oferta,_ := model.GetOferta(c, out.IdOft)
			if oferta.IdEmp == "none" {
				out.Status = "invalidUpload"
				berr := blobstore.Delete(c, file[0].BlobKey)
				model.Check(berr)
			} else {
				out.Status = "ok"
				if len(file) == 0 {
					out.Status = "invalidUpload"
					berr := blobstore.Delete(c, file[0].BlobKey)
					model.Check(berr)
				} else {
					var oldblobkey = oferta.BlobKey
					oferta.BlobKey = file[0].BlobKey
					out.IdOft = oferta.IdOft
					err = model.PutOferta(c, oferta)
					if err != nil {
						out.Status = "invalidUpload"
						berr := blobstore.Delete(c, file[0].BlobKey)
						model.Check(berr)
					}
					/* 
						Se borra el blob anterior, porque siempre crea uno nuevo
						No se necesita revisar el error
						Si es el blobkey = none no se borra por obvias razones
						Se genera una sesion nueva de upload en caso de que quieran
						cambiar la imágen en la misma pantalla. Esto es debido a que
						se utiliza un form estático con ajax
					*/
					if oldblobkey != "none" {
						blobstore.Delete(c, oldblobkey)
						UploadURL, err := blobstore.UploadURL(c, "/r/ofimgup", nil)
						out.UploadURL = UploadURL.String()
						if err != nil {
							out.Status = "uploadSessionError"
						}
					}
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(out)
	w.Write(b)
}
