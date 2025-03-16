package main

import (
	"fmt"
	"log"
	"net"
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
	forwardIP := os.Getenv("FORWARD_IP") == "true"
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
		if forwardIP {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			r.Header.Add("X-Real-IP", ip)
			r.Header.Add("X-Forwarded-For", ip)
		}
		proxy.ServeHTTP(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
