package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
)

func main() {
	n := negroni.Classic()
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/img", imageHandler)

	n.UseHandler(mux)

	n.Run(":3000")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	bgForm, _, _ := r.FormFile("background")
	imgForm, _, _ := r.FormFile("image")

	bg, _, _ := image.Decode(bgForm)
	img, _, _ := image.Decode(imgForm)

	hasher := sha256.New()
	bgFileName := hasher.Sum([]byte("bg"))
	bgFile, _ := os.Create(fmt.Sprintf("./storage/%x", bgFileName))
	defer bgFile.Close()
	gif.Encode(bgFile, bg, nil)

	imgFileName := hasher.Sum([]byte("img"))
	imgFile, _ := os.Create(fmt.Sprintf("./storage/%x", imgFileName))
	defer imgFile.Close()
	gif.Encode(imgFile, img, nil)
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
