package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hasanm95/caching-proxy-server/proxy"
)

func main() {
	port := flag.Int("port", 0, "Port to run the proxy server on")
	origin := flag.String("origin", "", "Origin server URL (e.g., http://example.com)")
	clearCache := flag.Bool("clear-cache", false, "Clear in-memory cache on start")

	flag.Parse()

	cache := proxy.NewCaheStore()

	if *clearCache {
		cache.Clear()
		log.Println("In-memory cache cleared")
	}

	if *origin == "" {
		log.Fatal("origin is required")
	}

	if *port == 0 {
		log.Fatal("port is required")
	}

	http.HandleFunc("/", proxy.HandleProxy(cache, *origin))

	log.Printf("Server will run at port %d and origin %s", *port, *origin)

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
