// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package micrositio

import (
	"appengine"
	"appengine/datastore"
	"crypto/sha1"
	"resize"
	"bytes"
	"strings"
	"strconv"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // import so we can read PNG files.
	"io"
	"net/http"
	"text/template"
	"sess"
	"model"
)

var (
	templates = template.Must(template.ParseFiles(
		"templates/error.html",
	))
)

type FormDataImage struct {
	Data	[]byte
	IdEmp	string
	IdImg	string
	Kind	string
	Name	string
	ErrName	string
	Desc	string
	ErrDesc	string
	Sizepx	int
	Sizepy	int
	Url		string
	ErrUrl	string
	Type	string
	Sp1		string
	Sp2		string
	Sp3		string
	Sp4		string
	Np1		int
	Np2		int
	Np3		int
	Np4		int
}
// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {
	http.HandleFunc("/mi", model.ErrorHandler(micrositio))
	http.HandleFunc("/logoup", model.ErrorHandler(upload))
	http.HandleFunc("/midata", model.ErrorHandler(modData))
	http.HandleFunc("/logosz", model.ErrorHandler(resizeLogo))
	http.HandleFunc("/simg", model.ErrorHandler(img))
}

var micrositioTpl = template.Must(template.ParseFiles("templates/micrositio.html")) //, "templates/login.html"))

func micrositio(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		if emp != nil {
			img := model.GetLogo(c, r.FormValue("IdEmp"))
			if(img == nil) {
				img = new(model.Image)
				img.IdEmp = emp.IdEmp
			}
			fd := imgToForm(*img)
			tc := make(map[string]interface{})
			tc["Sess"] =  s
			tc["Empresa"] = emp
			tc["FormData"] = fd
			micrositioTpl.Execute(w, tc)
		}
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

// upload is the HTTP handler for uploading images; it handles "/".
func upload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		imgo := model.GetLogo(c, r.FormValue("IdEmp"))
		if imgo == nil {
			imgo = new(model.Image)
			imgo.IdEmp = emp.IdEmp
		}
		fd := imgToForm(*imgo)
		tc := make(map[string]interface{})
		tc["Sess"] =  s
		tc["Empresa"] = emp
		tc["FormData"] = fd
		if r.Method != "POST" {
			// No upload; show the upload form.
			micrositio(w, r)
			return
		}

		idemp := r.FormValue("IdEmp")
		name := r.FormValue("Name")
		desc := r.FormValue("Desc")
		url := r.FormValue("Url")
		f, _, err := r.FormFile("image")
		model.Check(err)
		defer f.Close()

		// Grab the image data
		var buf bytes.Buffer
		io.Copy(&buf, f)
		i, _, err := image.Decode(&buf)
		if err != nil {
			if(r.FormValue("tipo")=="async") {
				//w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó el logotipo, formato no aceptado");
			} else {
				tc["Error"] = struct { Badformat string }{"badformat"}
				micrositioTpl.Execute(w, tc)
			}
			return
		}

		// Resize if too large, for more efficient moustachioing.
		// We aim for less than 1200 pixels in any dimension; if the
		// picture is larger than that, we squeeze it down to 600.
		const max = 800
		if b := i.Bounds(); b.Dx() > max || b.Dy() > max {
			// If it's gigantic, it's more efficient to downsample first
			// and then resize; resizing will smooth out the roughness.
			if b.Dx() > 2*max || b.Dy() > 2*max {
				w, h := max, max
				if b.Dx() > b.Dy() {
					h = b.Dy() * h / b.Dx()
				} else {
					w = b.Dx() * w / b.Dy()
				}
				i = resize.Resample(i, i.Bounds(), w, h)
				b = i.Bounds()
			}
			w, h := max/2, max/2
			if b.Dx() > b.Dy() {
				h = b.Dy() * h / b.Dx()
			} else {
				w = b.Dx() * w / b.Dy()
			}
			i = resize.Resize(i, i.Bounds(), w, h)
		} else {
			w, h := max, max
			if b.Dx() > b.Dy() {
				h = b.Dy() * h / b.Dx()
			} else {
				w = b.Dx() * w / b.Dy()
			}
			i = resize.Resize(i, i.Bounds(), w, h)
		}

		// Encode as a new JPEG image.
		buf.Reset()
		err = jpeg.Encode(&buf, i, nil)
		if err != nil {
			if(r.FormValue("tipo")=="async") {
				fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó el logotipo, formato no aceptado");
			} else {
				tc["Error"] = struct { Badencode string }{"badencode"}
				micrositioTpl.Execute(w, tc)
			}
			return
		}

		// Save the image under a unique key, a hash of the image.
		img := &model.Image{
			Data: buf.Bytes(), IdEmp: idemp, IdImg: model.RandId(12), 
			Kind: "EmpLogo", Name: name, Desc: desc, 
			Sizepx: 0, Sizepy: 0, Url: url, Type: "",
			Sp1: "", Sp2: "", Sp3: "", Sp4: "",
			Np1: 0, Np2: 0, Np3: 0, Np4: 0,
		}
		_, err = model.PutLogo(c, img)
		if err != nil {
			if(r.FormValue("tipo")=="async") {
				fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó el logotipo. Sistema en manetnimiento, intente en unos minutos");
			} else {
				tc["Error"] = struct { Cantsave string }{ "cantsave" }
				micrositioTpl.Execute(w, tc)
			}
			return
		}
		if(r.FormValue("tipo")=="async") {
			fmt.Fprintf(w, "<p></p>");
		} else {
			micrositio(w, r)
		}
		return
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

func modData(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		imgo := model.GetLogo(c, r.FormValue("IdEmp"))
		if(imgo == nil) {
			imgo = new(model.Image)
			imgo.IdEmp = emp.IdEmp
		}
		fd := imgToForm(*imgo)
		tc := make(map[string]interface{})
		tc["Sess"] =  s
		tc["Empresa"] = emp
		tc["FormData"] = fd
		if r.Method != "POST" {
			// No upload; show the upload form.
			micrositio(w, r)
			return
		}

		idemp := r.FormValue("IdEmp")
		name := r.FormValue("Name")
		desc := r.FormValue("Desc")
		url := r.FormValue("Url")

	//	key := datastore.NewKey(c, "EmpLogo", r.FormValue("id"), 0, nil)
	//	im := new(model.Image)
		// Save the image under a unique key, a hash of the image.
		imgo = &model.Image{
			Data: imgo.Data, IdEmp: idemp, 
			Kind: "EmpLogo", Name: name, Desc: desc, 
			Sizepx: 0, Sizepy: 0, Url: url, Type: "",
			Sp1: "", Sp2: "", Sp3: "", Sp4: "",
			Np1: 0, Np2: 0, Np3: 0, Np4: 0,
		}

		_, err := model.PutLogo(c, imgo)
		if err != nil {
			tc["Error"] = struct { Cantsave string }{ "cantsave" }
			micrositioTpl.Execute(w, tc)
			return
		}

		micrositio(w, r)
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

// upload is the HTTP handler for uploading images; it handles "/".
func resizeLogo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if s, ok := sess.IsSess(w, r, c); ok {
		emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
		imgo := model.GetLogo(c, r.FormValue("IdEmp"))
		if(imgo == nil) {
			imgo = new(model.Image)
			imgo.IdEmp = emp.IdEmp
		}
		sf, _ := strconv.Atoi(r.FormValue("s"))
		fd := imgToForm(*imgo)
		tc := make(map[string]interface{})
		tc["Sess"] =  s
		tc["Empresa"] = emp
		tc["FormData"] = fd
		if r.Method != "POST" {
			// No upload; show the upload form.
			micrositio(w, r)
			return
		}
		i, _, err := image.Decode(bytes.NewBuffer(imgo.Data))
		if err != nil {
			tc["Error"] = struct { Badformat string }{"badformat"}
			micrositioTpl.Execute(w, tc)
			return
		}

		// Resize if too large, for more efficient moustachioing.
		// We aim for less than 1200 pixels in any dimension; if the
		// picture is larger than that, we squeeze it down to 600.
		const max = 800
		if b := i.Bounds(); b.Dx() > max || b.Dy() > max {
			// If it's gigantic, it's more efficient to downsample first
			// and then resize; resizing will smooth out the roughness.
			if b.Dx() > 2*max || b.Dy() > 2*max {
				w, h := max, max
				if b.Dx() > b.Dy() {
					h = b.Dy() * h / b.Dx()
				} else {
					w = b.Dx() * w / b.Dy()
				}
				i = resize.Resample(i, i.Bounds(), w, h)
				b = i.Bounds()
			}
			w, h := max/2, max/2
			if b.Dx() > b.Dy() {
				h = b.Dy() * h / b.Dx()
			} else {
				w = b.Dx() * w / b.Dy()
			}
			i = resize.Resize(i, i.Bounds(), w, h)
		} else {
			h := b.Dy()
			w := b.Dx()
			if(sf > 0 && sf <= 2) {
				h = h * sf
				w = w * sf
			} else if (sf < 0 && sf > -1){
				sf = (sf*2)+sf
				h = h * (1/sf)
				w = w * (1/sf)
			}
			i = resize.Resize(i, i.Bounds(), w, h)
		}

		// Encode as a new JPEG image.
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, i, nil)
		if err != nil {
			tc["Error"] = struct { Badencode string }{"badencode"}
			micrositioTpl.Execute(w, tc)
			return
		}

		// Save the image under a unique key, a hash of the image.
		img := &model.Image{
			Data: buf.Bytes(), IdEmp: emp.IdEmp, IdImg: imgo.IdImg, 
			Kind: "EmpLogo", Name: imgo.Name, Desc: imgo.Desc, 
			Sizepx: 0, Sizepy: 0, Url: imgo.Url, Type: "",
			Sp1: "", Sp2: "", Sp3: "", Sp4: "",
			Np1: 0, Np2: 0, Np3: 0, Np4: 0,
		}
		_, err = model.PutLogo(c, img)
		if err != nil {
			tc["Error"] = struct { Cantsave string }{ "cantsave" }
			micrositioTpl.Execute(w, tc)
			return
		}

		micrositio(w, r)
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

// keyOf returns (part of) the SHA-1 hash of the data, as a hex string.
func keyOf(data []byte) string {
	sha := sha1.New()
	sha.Write(data)
	return fmt.Sprintf("%x", string(sha.Sum(nil))[0:8])
}

// img is the HTTP handler for displaying images and painting moustaches;
// it handles "/img".
func img(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "EmpLogo", r.FormValue("id"), 0, nil)
	im := new(model.Image)
	err := datastore.Get(c, key, im)
	model.Check(err)

	m, _, err := image.Decode(bytes.NewBuffer(im.Data))
	model.Check(err)

	// Ejemplo de cómo mejorar la conversión de datos 
	//get := func(n string) int { // helper closure
//		i, _ := strconv.Atoi(r.FormValue(n))
//		return i
//	}
//	s := get("s")

	w.Header().Set("Content-type", "image/jpeg")
	jpeg.Encode(w, m, nil)
}


func delimg(w http.ResponseWriter, r *http.Request) {
}

// errorHandler wraps the argument handler with an error-catcher that
// returns a 500 HTTP error if the request fails (calls check with err non-nil).
/*
func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if _, ok := recover().(error); ok {
				w.WriteHeader(http.StatusInternalServerError)
				c := appengine.NewContext(r)
				if s, ok := sess.IsSess(w, r, c); ok {
					emp := model.GetEmpresa(c, r.FormValue("IdEmp"))
					if emp != nil {
						tc := make(map[string]interface{})
						tc["Sess"] =  s
						tc["Empresa"] = emp
						tc["Error"] = struct { string }{""}
						micrositioTpl.Execute(w, tc)
					}
				} else {
					http.Redirect(w, r, "/registro", http.StatusFound)
				}
			}
		}()
		fn(w, r)
	}
}
*/

func imgForm(w http.ResponseWriter, r *http.Request, s sess.Sess, valida bool, tpl *template.Template) (FormDataImage, bool){
	fd := FormDataImage {
		IdEmp: strings.TrimSpace(r.FormValue("IdEmp")),
		Name: strings.TrimSpace(r.FormValue("Name")),
		ErrName: "",
		Url: strings.TrimSpace(r.FormValue("Url")),
		ErrUrl: "",
		Desc: strings.TrimSpace(r.FormValue("Desc")),
		ErrDesc: "",
	}
	if valida {
		var ef bool
		ef = false
		if fd.Name != "" && !model.ValidName.MatchString(fd.Name) {
			fd.ErrName = "invalid"
			ef = true
		}
		if fd.Url != "" && !model.ValidUrl.MatchString(fd.Url) {
			fd.ErrUrl = "invalid"
			ef = true
		}
		if fd.Desc != "" && !model.ValidSimpleText.MatchString(fd.Desc) {
			fd.ErrDesc = "invalid"
			ef = true
		}
		if ef {
			tc := make(map[string]interface{})
			tc["Sess"] = s
			tc["FormData"] = fd
			tpl.Execute(w, tc)
			return fd, false
		}
	}
	return fd, true
}
func imgFill(r *http.Request, img *model.Image) {
	img.Name=		strings.TrimSpace(r.FormValue("Name"))
	img.Desc=		strings.TrimSpace(r.FormValue("Desc"))
	img.Url=		strings.TrimSpace(r.FormValue("Url"))
}

func imgToForm(e model.Image) FormDataImage {
	fd := FormDataImage {
		IdEmp:		e.IdEmp,
		IdImg:		e.IdImg,
		Kind:		e.Kind,
		Name:		e.Name,
		Desc:		e.Desc,
		Sizepx:		e.Sizepx,
		Sizepy:		e.Sizepy,
		Url:		e.Url,
		Type:		e.Type,
		Sp1:		e.Sp1,
		Sp2:		e.Sp2,
		Sp3:		e.Sp3,
		Sp4:		e.Sp4,
		Np1:		e.Np1,
		Np2:		e.Np2,
		Np3:		e.Np3,
		Np4:		e.Np4,
	}
	return fd
}
