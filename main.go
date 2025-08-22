package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {
	port := flag.Int("port", 0, "Port to run the proxy server on")
	origin := flag.String("origin", "", "Origin server URL (e.g., http://example.com)")
	clearCache := flag.Bool("clear-cache", false, "Clear in-memory cache on start")

	flag.Parse()

	cache := NewCaheStore()

	if *clearCache {
		cache.clear()
		log.Println("In-memory cache cleared")
	}

	if *origin == "" {
		log.Fatal("origin is required")
	}

	if *port == 0 {
		log.Fatal("port is required")
	}

	http.HandleFunc("/", handleProxy(cache, *origin))

	log.Printf("Server will run at port %d and origin %s", *port, *origin)

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}

// Cached response in memory
type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type CacheStore struct {
	sync.RWMutex
	Store map[string]CachedResponse
}

func (c *CacheStore) set(key string, resp CachedResponse) {
	c.Lock()
	defer c.Unlock()
	c.Store[key] = resp
}

func (c *CacheStore) get(key string) (CachedResponse, bool) {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.Store[key]
	return v, ok
}

func (c *CacheStore) clear() {
	c.Lock()
	defer c.Unlock()
	c.Store = make(map[string]CachedResponse)
}

func NewCaheStore() *CacheStore {
	return &CacheStore{Store: make(map[string]CachedResponse)}
}

func cloneHeader(h http.Header) http.Header {
	nh := make(http.Header)
	for k, v := range h {
		v2 := make([]string, len(v))
		copy(v2, v)
		nh[k] = v2
	}
	return nh
}

func handleProxy(cache *CacheStore, origin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET is cached/proxied for now", http.StatusMethodNotAllowed)
			return
		}
		orginUrl := strings.TrimRight(origin, "/") + r.URL.RequestURI()

		// cache
		cacheKey := r.Method + ":" + orginUrl
		if cachedData, ok := cache.get(cacheKey); ok {
			for k, vals := range cachedData.Headers {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(cachedData.StatusCode)
			w.Write(cachedData.Body)
			return
		}

		// request to origin
		req, err := http.NewRequest(r.Method, orginUrl, nil)

		if err != nil {
			http.Error(w, "Origin request error", http.StatusBadGateway)
			return
		}

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			http.Error(w, "Origin unreachable: "+err.Error(), http.StatusBadGateway)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, "Origin body read error: "+err.Error(), http.StatusBadGateway)
			return
		}

		// forward headers
		for k, vals := range resp.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}

		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

		// Save to cache
		cache.set(cacheKey, CachedResponse{
			StatusCode: resp.StatusCode,
			Headers:    cloneHeader(resp.Header),
			Body:       body,
		})
	}
}
