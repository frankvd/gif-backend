package main

import (
	"encoding/json"
	"os"
)

// Config struct
type Config struct {
	HMACSecret   string `json:"hmac_secret"`
	FrontendHost string `json:"frontend_host"`
	RedisHost    string `json:"redis_host"`
}

// Read json file into Config struct
func readConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)

	if err != nil {
		panic(err)
	}

	return conf
}
