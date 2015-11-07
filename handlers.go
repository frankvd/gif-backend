package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"net/http"
	"os"
	"strconv"
)

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

	x, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		x = 0
	}
	y, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		y = 0
	}
	pt.X = x
	pt.Y = y
	draw.Draw(m, m.Bounds(), src, *pt, draw.Over)

	buffer := new(bytes.Buffer)
	gif.Encode(buffer, m, nil)

	fmt.Fprint(w, base64.StdEncoding.EncodeToString(buffer.Bytes()))
}
