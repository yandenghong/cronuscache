package main

import (
	"fmt"
	"log"
	"net/http"
	"cronuscache/core"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	src.NewGroup("scores", 2<<10, src.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := src.NewHTTPPool(addr)
	log.Println("cronuscache is running at", addr)
	http.ListenAndServe(addr, peers)
}
