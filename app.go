package main

import (
	"crypto/sha256"
	"hash"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/garyburd/redigo/redis"
)

var hasher hash.Hash
var rds redis.Conn
var config Config

func main() {
	config = readConfig("./config.json")
	rds, _ = redis.Dial("tcp", config.RedisHost)
	hasher = sha256.New()
	n := negroni.Classic()
	// Register middleware
	n.UseFunc(corsMiddleware)
	n.UseFunc(authMiddleware)
	n.UseFunc(rateLimitMiddleware)
	mux := http.NewServeMux()
	// Register rotues
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/img", imageHandler)
	n.UseHandler(mux)

	// Run app
	n.Run(":3000")
}
