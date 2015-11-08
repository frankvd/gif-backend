package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
)

// Returns the HMAC key set in the config
func hmacKeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}
	return []byte(config.HMACSecret), nil
}

// Validates the JWT token in the Authorization header
func authMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := jwt.ParseFromRequest(r, hmacKeyFunc)

	if err != nil {
		panic(err)
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	next(w, r)
}

// Sets correct CORS headers
func corsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", config.FrontendHost)
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	if r.Method == "OPTIONS" {
		return
	}

	next(w, r)
}

// Rate limit middleware
func rateLimitMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := jwt.ParseFromRequest(r, hmacKeyFunc)
	if err != nil {
		panic(err)
	}

	if !rateLimit(token.Claims["iss"].(string)) {
		w.WriteHeader(429)
		return
	}

	next(w, r)
}

// Rate limit api call
func rateLimit(username string) bool {
	limitKey := username + ":" + time.Now().Format("2006-01-02")
	reply, err := rds.Do("GET", limitKey)

	if err != nil {
		panic(err)
	}

	current, _ := redis.Int(reply, nil)

	if current > 5 {
		return false
	}
	rds.Do("INCR", limitKey)
	rds.Do("EXPIRE", limitKey, 86400)
	return true
}
