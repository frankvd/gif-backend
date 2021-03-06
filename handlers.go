package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

func getFileName(imgName string) string {
	dir := "./storage/"
	imgFileName := hasher.Sum([]byte(imgName))

	return fmt.Sprintf(dir+"%x", imgFileName)
}

func sendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

// Upload handler
// Stores the POSTed GIF and overlay
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := jwt.ParseFromRequest(r, hmacKeyFunc)
	issuer := token.Claims["iss"].(string)

	gfForm, _, err1 := r.FormFile("gif")
	overlayForm, _, err2 := r.FormFile("overlay")

	if err1 != nil || err2 != nil {
		sendResponse(w, 400, map[string]string{
			"error": "Error while reading images",
		})
		return
	}

	gf, err1 := gif.DecodeAll(gfForm)
	overlay, _, err2 := image.Decode(overlayForm)

	if err1 != nil || err2 != nil {
		sendResponse(w, 400, map[string]string{
			"error": "Error while decoding images",
		})
		return
	}

	gfFile, err1 := os.Create(getFileName("gif:" + issuer))
	defer gfFile.Close()
	gif.EncodeAll(gfFile, gf)

	overlayFile, err2 := os.Create(getFileName("overlay:" + issuer))
	defer overlayFile.Close()
	png.Encode(overlayFile, overlay)

	if err1 != nil || err2 != nil {
		sendResponse(w, 400, map[string]string{
			"error": "Error while saving images",
		})
		return
	}
}

// Image Handler
// Reads the stored GIF & overlay for the authenticated user and returns a new GIF with the overlay copied onto each frame of the input GIF
func imageHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := jwt.ParseFromRequest(r, hmacKeyFunc)
	issuer := token.Claims["iss"].(string)

	gfFile, err1 := os.Open(getFileName("gif:" + issuer))
	overlayFile, err2 := os.Open(getFileName("overlay:" + issuer))

	if err1 != nil || err2 != nil {
		sendResponse(w, 400, map[string]string{
			"error": "Error while reading images",
		})
		return
	}

	// Decode images
	dst, err1 := gif.DecodeAll(gfFile)
	src, _, err2 := image.Decode(overlayFile)

	if err1 != nil || err2 != nil {
		sendResponse(w, 400, map[string]string{
			"error": "Error while decoding images",
		})
		return
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
