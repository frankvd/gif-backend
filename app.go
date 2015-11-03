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

var hasher = sha256.New()

func main() {
	n := negroni.Classic()
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/img", imageHandler)

	n.UseHandler(mux)

	n.Run(":3000")
}

func getFileName(imgName string) string {
	dir := "./storage/"
	imgFileName := hasher.Sum([]byte(imgName))

	return fmt.Sprintf(dir+"%x", imgFileName)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	bgForm, _, _ := r.FormFile("background")
	imgForm, _, _ := r.FormFile("image")

	bg, _, _ := image.Decode(bgForm)
	img, _, _ := image.Decode(imgForm)

	bgFile, _ := os.Create(getFileName("bg"))
	defer bgFile.Close()
	gif.Encode(bgFile, bg, nil)

	imgFile, _ := os.Create(getFileName("img"))
	defer imgFile.Close()
	gif.Encode(imgFile, img, nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	bgFile, _ := os.Open(getFileName("bg"))
	imgFile, _ := os.Open(getFileName("img"))

	dst, _, _ := image.Decode(bgFile)
	src, _, _ := image.Decode(imgFile)
	m := image.NewRGBA(dst.Bounds())
	pt := new(image.Point)
	draw.Draw(m, m.Bounds(), dst, *pt, draw.Src)
	draw.Draw(m, m.Bounds(), src, *pt, draw.Over)

	w.Header().Set("Content-Type", "image/gif")

	gif.Encode(w, m, nil)
}
