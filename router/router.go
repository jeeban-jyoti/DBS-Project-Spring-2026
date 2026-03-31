package router

import (
	"fmt"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/authentication"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/student"
)

var publicRoutes = map[string]http.HandlerFunc{
	"/api/v1/login":  authentication.Login,
	"/api/v1/logout": authentication.Logout,
}

var protectedRoutes = map[string]http.HandlerFunc{
	"/api/v1/test": test,

	"/api/v1/fetchBooks": student.FetchAllBooks,
}

// test api
func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from protected route")
}

func Route(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	if handler, ok := publicRoutes[r.URL.Path]; ok {
		handler(w, r)
		return
	}

	if handler, ok := protectedRoutes[r.URL.Path]; ok {
		authentication.RequireAuth(handler)(w, r)
		return
	}

	http.NotFound(w, r)
}
