package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	host := os.Getenv("PROXY_HOST")
	if host == "" {
		panic("PROXY_HOST is not set")
	}
	schema := os.Getenv("PROXY_SCHEMA")
	if schema == "" {
		schema = "https"
	}
	u := schema + "://" + host
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	// initialize a reverse proxy and pass the actual backend server url here
	proxy := httputil.NewSingleHostReverseProxy(url)

	// handle all requests to your server using the proxy
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = host
		proxy.ServeHTTP(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
