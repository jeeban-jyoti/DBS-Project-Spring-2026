package router

import (
	"fmt"
	"net/http"

	auth "github.com/jeeban-jyoti/DSB-Project-Spring-2026/authentication"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/student"
)

var publicRoutes = map[string]http.HandlerFunc{
	"/api/v1/login":  auth.Login,
	"/api/v1/logout": auth.Logout,
}

var protectedRoutes = map[string]http.HandlerFunc{
	"/api/v1/test": test,

	"/api/v1/fetchBooks":      student.FetchAllBooks,
	"api/v1/addToCart":        auth.RequireRole(auth.RoleStudent)(student.AddToCart),
	"api/v1/removeFromCart":   auth.RequireRole(auth.RoleStudent)(student.RemoveFromCart),
	"api/v1/placeBuyOrder":    auth.RequireRole(auth.RoleStudent)(student.PlaceBuyOrder),
	"api/v1/placeBorrowOrder": auth.RequireRole(auth.RoleStudent)(student.PlaceBorrowOrder),
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
		auth.RequireAuth(handler)(w, r)
		return
	}

	http.NotFound(w, r)
}
