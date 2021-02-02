package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
)

func serviceNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	fmt.Fprintf(w, "{ \"error\": \"service not found\" }")
}

func redirectRequest(w http.ResponseWriter, name string, path string, r *http.Request) {
    var url = fmt.Sprintf("http://%s/%s", name, path)

    http.Redirect(w, r, url, http.StatusFound)
}

func handler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path[1:])
	if err != nil {
        log.Fatal(err)
    }

    urlSegments := strings.SplitAfterN(u.Path, "/", 2)

	var serviceName string = urlSegments[0]

    if len(urlSegments) > 1 {
	    servicePath := urlSegments[0]
	    serviceName = servicePath[:len(servicePath)-1]
    }

    if len(serviceName) < 1 {
    	serviceNotFound(w)
    } else {
    	var path string = ""

    	if len(urlSegments) > 1 {
    		path = urlSegments[1]
    	}

    	redirectRequest(w, serviceName, path, r)
    }
}

func main() {
	var port int = 8080

    http.HandleFunc("/", handler)

	fmt.Printf("listening on port %d\n", port)

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}