package main

import (
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/codegangsta/negroni"
)

func main() {
	n := negroni.Classic()
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", imageHandler)

	n.UseHandler(mux)

	n.Run(":3000")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	overlay, _, _ := r.FormFile("overlay")
	gf, _, err := r.FormFile("gif")
	if err != nil {
		panic(err)
	}

	dst, _, _ := image.Decode(overlay)
	src, _, _ := image.Decode(gf)
	m := image.NewRGBA(dst.Bounds())
	pt := new(image.Point)
	draw.Draw(m, m.Bounds(), dst, *pt, draw.Src)
	draw.Draw(m, m.Bounds(), src, *pt, draw.Over)

	w.Header().Set("Content-Type", "image/gif")

	gif.Encode(w, m, nil)
}
