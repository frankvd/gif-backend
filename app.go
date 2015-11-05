package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
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

func hmacKeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}
	return []byte("super_secret_key"), nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	token, err := jwt.ParseFromRequest(r, hmacKeyFunc)

	if err != nil || !token.Valid {
		fmt.Printf("%v", token.Valid)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	if r.Method == "OPTIONS" {
		return
	}

	token, err := jwt.ParseFromRequest(r, hmacKeyFunc)

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
