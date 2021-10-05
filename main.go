package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func serviceNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, err := fmt.Fprintf(w, "{ \"error\": \"service not found\" }")
	if err != nil {
		return
	}
}

func redirectRequest(protocol string, w http.ResponseWriter, host string, path string, r *http.Request) {
	var newUrl = fmt.Sprintf("%s%s%s?%s", protocol, host, path, r.URL.RawQuery)

	remote, err := url.Parse(newUrl)
	if err != nil {
		log.Print(err)
		serviceNotFound(w)
		return
	}

	director := func(req *http.Request) {
		req.Host = host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = host
		req.URL.Path = path

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "simple-proxy")
		}
	}

	proxy := &httputil.ReverseProxy{Director: director}

	log.Print(fmt.Sprintf("redirecting to: %s\n", newUrl))

	proxy.ServeHTTP(w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("APP_HOST")
	protocol := os.Getenv("APP_PROTOCOL")
	path := r.URL.Path
	redirectRequest(protocol, w, host, path, r)
}

func main() {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal("Error converting port to int")
	}

	http.HandleFunc("/", handler)

	fmt.Printf("listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
