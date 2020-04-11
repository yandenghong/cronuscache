package main

import (
	"cronuscache/core"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *core.Group {
	return core.NewGroup("scores", 2<<10, core.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// startCacheServer start the cache server, create HTTPPool,
// add node information, register with groups, and start the HTTP service (3 ports, 8309/8310/8311)
func startCacheServer(addr string, addrs []string, g *core.Group) {
	nodes := core.NewHTTPPool(addr)
	nodes.Set(addrs...)
	g.RegisterNodes(nodes)
	log.Println("cronuscache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nodes))
}

// startAPIServer start an API service (port 9999) that interacts with the user
func startAPIServer(apiAddr string, g *core.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := g.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8309, "CronusCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8309: "http://localhost:8309",
		8310: "http://localhost:8310",
		8311: "http://localhost:8311",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	g := createGroup()
	if api {
		go startAPIServer(apiAddr, g)
	}
	startCacheServer(addrMap[port], []string(addrs), g)
}
