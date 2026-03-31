package router

import (
	"fmt"
	"net/http"
)

var routes = map[string]http.HandlerFunc{
	"/api/v1/test": test,
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from mapped side here")
}

func Route(w http.ResponseWriter, r *http.Request) {
	if handler, ok := routes[r.URL.Path]; ok {
		handler(w, r)
		return
	}
	http.NotFound(w, r)
}
