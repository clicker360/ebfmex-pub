// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On App Engine, the framework sets up main; we should be a different package.
package micrositio

import (
	"appengine"
	//"appengine/datastore"
	//"crypto/sha1"
	"resize"
	"bytes"
	//"strings"
	"strconv"
	//"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // import so we can read PNG files.
	//"io"
	"net/http"
	//"text/template"
	//"sess"
	"model"
)

func init() {
	http.HandleFunc("/spic", rslogo)
}

// Resize only if picture from EmpLogo is more than 80 pix width
// If resize, save image to entity king ShortLogo and stream
// If no resize is necesary, save to entity ShortLogo and stream
func rslogo(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.Method == "GET" {
		sf, _ := strconv.Atoi(r.FormValue("s"))
		// Check for shortlogo
		var simg model.Image
		simg.Kind = "ShortLogo"
		if(r.FormValue("IdEmp")!="") {
			simg.IdEmp = r.FormValue("IdEmp");
		} else {
			return
		}
		shortlogo, _ := model.GetShortLogo(c, r.FormValue("IdEmp"))

		// Stream short logo if already exists
		if(shortlogo != nil) {
			m, _, err := image.Decode(bytes.NewBuffer(shortlogo.Data))
			check(err)
			w.Header().Set("Content-type", "image/jpeg")
			jpeg.Encode(w, m, nil)
			return
		}

		// Process biglogo if shortlogo doesn't exists 
		// Save and Stream new shortlogo
		biglogo := model.GetLogo(c, r.FormValue("IdEmp"))
		if biglogo == nil {
			// No such imageb
			return
		}

		i, _, err := image.Decode(bytes.NewBuffer(biglogo.Data))
		if err != nil {
			// ERR_BADFORMAT
			return
		}

		const max = 80
		if(sf==0) {
			// We aim for less than 80 pixels in any dimension.
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
		} else {
			// We aim for a resize by ratio.
			b := i.Bounds()
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
			// ERR_FAIL_ENCODE
			return
		}

		// Save the image under a unique key, a hash of the image.
		simg.IdImg = model.RandId(12)
		simg.Data = buf.Bytes()
		simg.Name = biglogo.Name
		simg.Desc = biglogo.Desc
		simg.Sizepx = max
		simg.Sizepy = 0
		simg.Url = biglogo.Url
		simg.Type = "jpeg"
		simg.Sp1 = biglogo.Sp1
		simg.Sp2 = biglogo.Sp2
		simg.Sp3 = biglogo.Sp3
		simg.Sp4 = biglogo.Sp4
		simg.Np1 = biglogo.Np1
		simg.Np2 = biglogo.Np2
		simg.Np3 = biglogo.Np3
		simg.Np4 = biglogo.Np4

		_, err = model.PutLogo(c, &simg)
		if err != nil {
			// ERR_DATASTORE
			// Don't return, stream image either way
		}
		m, _, err := image.Decode(bytes.NewBuffer(simg.Data))
		check(err)
		w.Header().Set("Content-type", "image/jpeg")
		jpeg.Encode(w, m, nil)
	}
}
