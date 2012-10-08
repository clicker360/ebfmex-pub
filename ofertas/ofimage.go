// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package oferta

import (
	"appengine"
	"resize"
	"bytes"
	//"strconv"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // import so we can read PNG files.
	"io"
	"net/http"
	"sess"
	"model"
)

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {
	http.HandleFunc("/ofimgup", model.ErrorHandler(upload))
	http.HandleFunc("/ofimg", model.ErrorHandler(img))
}

// upload is the HTTP handler for uploading images; it handles "/".
func upload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if _, ok := sess.IsSess(w, r, c); ok {
		oferta := model.GetOferta(c, r.FormValue("IdOft"))
		if oferta.IdEmp == "none" {
			// No se pude agregar imagen si la oferta es un cascaron
			return
		}
		if r.Method != "POST" {
			return
		}

		f, _, err := r.FormFile("image")
		model.Check(err)
		defer f.Close()

		// Grab the image data
		var buf bytes.Buffer
		io.Copy(&buf, f)
		i, _, err := image.Decode(&buf)
		if err != nil {
			fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó la imagen, formato no aceptado");
			return
		}

		// Resize if too large, for more efficient moustachioing.
		// We aim for less than 1200 pixels in any dimension; if the
		// picture is larger than that, we squeeze it down to 600.
		const max = 600
		// We aim for less than max pixels in any dimension.
		if b := i.Bounds(); b.Dx() > max || b.Dy() > max {
			// If it's gigantic, it's more efficient to downsample first
			// and then resize; resizing will smooth out the roughness.
			if b.Dx() > 2*max || b.Dy() > 2*max {
				w, h := max*2, max*2
				if b.Dx() > b.Dy() {
					h = b.Dy() * h / b.Dx()
				} else {
					w = b.Dx() * w / b.Dy()
				}
				i = resize.Resample(i, i.Bounds(), w, h)
				b = i.Bounds()
			}
			w, h := max, max
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
		b := i.Bounds()
		err = jpeg.Encode(&buf, i, nil)
		if err != nil {
			fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó la imagen, formato no aceptado");
			return
		}

		// Save the image under a unique key, a hash of the image.
		oferta.Image = buf.Bytes()
		oferta.Sizepx = b.Dx()
		oferta.Sizepy = b.Dy()

		err = model.PutOferta(c, oferta)
		if err != nil {
			fmt.Fprintf(w, "<p>'%s'</p>", "No se actualizó la imagen. Sistema en matenimiento, intente en unos minutos");
			return
		}

		fmt.Fprintf(w, "<p></p>");
		return
	} else {
		http.Redirect(w, r, "/registro", http.StatusFound)
	}
}

// img is the HTTP handler for displaying images;
// it handles "/simg".
func img(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	oferta := model.GetOferta(c, r.FormValue("id"))
	if oferta.IdEmp == "none" {
		// No se pude agregar imagen si la oferta es un cascaron
		return
	}
	m, _, err := image.Decode(bytes.NewBuffer(oferta.Image))
	model.Check(err)

	w.Header().Set("Content-type", "image/jpeg")
	jpeg.Encode(w, m, nil)
}
