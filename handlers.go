package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
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
	gfForm, _, err := r.FormFile("gif")
	overlayForm, _, err := r.FormFile("overlay")

	if err != nil {
		fmt.Fprint(w, "error reading images")
		panic(err)
	}

	gf, _ := gif.DecodeAll(gfForm)
	overlay, _, _ := image.Decode(overlayForm)

	gfFile, _ := os.Create(getFileName("gif"))
	defer gfFile.Close()
	gif.EncodeAll(gfFile, gf)

	overlayFile, _ := os.Create(getFileName("overlay"))
	defer overlayFile.Close()
	png.Encode(overlayFile, overlay)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	gfFile, err := os.Open(getFileName("gif"))
	overlayFile, err := os.Open(getFileName("overlay"))

	if err != nil {
		fmt.Fprint(w, "error reading images")
		panic(err)
	}

	// Decode images
	dst, err := gif.DecodeAll(gfFile)
	src, _, err := image.Decode(overlayFile)

	if err != nil {
		fmt.Fprint(w, "error decoding images")
		panic(err)
	}

	// Create location for the draw
	pt := new(image.Point)
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
	// Draw the overlay over each frame of the GIF
	for _, frame := range dst.Image {
		draw.Draw(frame, frame.Bounds(), src, *pt, draw.Over)
	}

	// Encode image to base64
	buffer := new(bytes.Buffer)
	gif.EncodeAll(buffer, dst)
	fmt.Fprint(w, base64.StdEncoding.EncodeToString(buffer.Bytes()))
}
