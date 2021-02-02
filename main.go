package main

import (
    "fmt"
    "log"
    "net/url"
    "net/http"
    "net/http/httputil"
    "strings"
)

func serviceNotFound(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)

    fmt.Fprintf(w, "{ \"error\": \"service not found\" }")
}

func redirectRequest(w http.ResponseWriter, name string, path string, r *http.Request) {
    var newUrl = fmt.Sprintf("http://%s/%s?%s", name, path, r.URL.RawQuery)

    remote, err := url.Parse(newUrl)
    if err != nil {
        log.Print(err)
        serviceNotFound(w)
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)

    log.Print(fmt.Sprintf("redirecting to: %s\n", newUrl))

    w.Header().Set("X-Ben", "Rad")
    proxy.ServeHTTP(w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
    urlSegments := strings.SplitAfterN(r.URL.Path[1:], "/", 2)

    var serviceName string = urlSegments[0]

    if len(serviceName) < 1 {
        serviceNotFound(w)
        return
    }

    var path string = ""

    if len(urlSegments) > 1 {
        servicePath := urlSegments[0]
        serviceName = servicePath[:len(servicePath)-1]
        path = urlSegments[1]
    }

    redirectRequest(w, serviceName, path, r)
}

func main() {
    var port int = 8080

    http.HandleFunc("/", handler)

    fmt.Printf("listening on port %d\n", port)

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}