package router

import (
	"fmt"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/admin"
	auth "github.com/jeeban-jyoti/DSB-Project-Spring-2026/authentication"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/orders"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/student"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/superadmin"
	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/tickets"
)

var publicRoutes = map[string]http.HandlerFunc{
	"/api/v1/login":  auth.Login,
	"/api/v1/logout": auth.Logout,
}

var protectedRoutes = map[string]http.HandlerFunc{
	"/api/v1/test": test,

	"/api/v1/changePassword": auth.ChangePassword,

	"/api/v1/fetchBooks":     student.FetchAllBooks,
	"/api/v1/fetchBook":      student.FetchBook,
	"/api/v1/cart":           auth.RequireRole(auth.RoleStudent)(student.ShowCart),
	"/api/v1/addToCart":      auth.RequireRole(auth.RoleStudent)(student.AddToCart),
	"/api/v1/removeFromCart": auth.RequireRole(auth.RoleStudent)(student.RemoveFromCart),

	"/api/v1/placeBuyOrder":     auth.RequireRole(auth.RoleStudent)(orders.PlaceBuyOrder),
	"/api/v1/placeBorrowOrder":  auth.RequireRole(auth.RoleStudent)(orders.PlaceBorrowOrder),
	"/api/v1/showMyOrders":      auth.RequireRole(auth.RoleStudent)(orders.ShowOrders),
	"/api/v1/showAllOrders":     auth.RequireRole(auth.RoleSupport)(orders.ShowAllOrders),
	"/api/v1/cancelOrder":       auth.RequireRole(auth.RoleStudent)(orders.GenerateOrderCancellation),
	"/api/v1/changeOrderStatus": auth.RequireRole(auth.RoleSupport)(orders.ChangeOrderStatus),
	"/api/v1/bookReturn":        auth.RequireRole(auth.RoleSupport)(orders.ReturnBorrowedBook),

	"/api/v1/generateTicket":             auth.RequireRole(auth.RoleStudent, auth.RoleSupport)(tickets.GenerateNewTicket),
	"/api/v1/viewMyTickets":              auth.RequireRole(auth.RoleStudent, auth.RoleSupport)(tickets.ShowGeneratedTickets),
	"/api/v1/viewALLTickets":             auth.RequireRole(auth.RoleAdmin)(tickets.ShowALLTickets),
	"/api/v1/viewNewTickets":             auth.RequireRole(auth.RoleSupport)(tickets.ShowNewTickets),
	"/api/v1/handleNewTicket":            auth.RequireRole(auth.RoleSupport)(tickets.AssignTicket),
	"/api/v1/changeAssignedTicketStatus": auth.RequireRole(auth.RoleAdmin)(tickets.ChangeTicketStatus),

	"/api/v1/addAdmin":           auth.RequireRole(auth.RoleSuperAdmin)(superadmin.AddAdmin),
	"/api/v1/removeAdmin":        auth.RequireRole(auth.RoleSuperAdmin)(superadmin.RemoveAdmin),
	"/api/v1/addSupportStaff":    auth.RequireRole(auth.RoleSuperAdmin)(superadmin.AddSupportStaff),
	"/api/v1/removeSupportStaff": auth.RequireRole(auth.RoleSuperAdmin)(superadmin.RemoveSupportStaff),

	"/api/v1/addUniversity":    auth.RequireRole(auth.RoleAdmin)(admin.AddUniversity),
	"/api/v1/removeUniversity": auth.RequireRole(auth.RoleAdmin)(admin.RemoveUniversity),
	"/api/v1/updateUniversity": auth.RequireRole(auth.RoleAdmin)(admin.UpdateUniversity),

	"/api/v1/addBook":    auth.RequireRole(auth.RoleAdmin)(admin.AddBook),
	"/api/v1/removeBook": auth.RequireRole(auth.RoleAdmin)(admin.RemoveBook),

	"/api/v1/addDepartment":    auth.RequireRole(auth.RoleAdmin)(admin.AddDepartment),
	"/api/v1/removeDepartment": auth.RequireRole(auth.RoleAdmin)(admin.RemoveDepartment),

	"/api/v1/addCourse":    auth.RequireRole(auth.RoleAdmin)(admin.AddCourse),
	"/api/v1/removeCourse": auth.RequireRole(auth.RoleAdmin)(admin.RemoveCourse),

	"/api/v1/addInstructor":    auth.RequireRole(auth.RoleAdmin)(admin.AddInstructor),
	"/api/v1/removeInstructor": auth.RequireRole(auth.RoleAdmin)(admin.RemoveInstructor),
	"/api/v1/updateInstructor": auth.RequireRole(auth.RoleAdmin)(admin.UpdateInstructor),

	"/api/v1/addStudent":    auth.RequireRole(auth.RoleAdmin)(admin.AddStudent),
	"/api/v1/removeStudent": auth.RequireRole(auth.RoleAdmin)(admin.RemoveStudent),
	"/api/v1/updateStudent": auth.RequireRole(auth.RoleAdmin)(admin.UpdateStudent),

	"/api/v1/addSemester":    auth.RequireRole(auth.RoleAdmin)(admin.AddSemester),
	"/api/v1/removeSemester": auth.RequireRole(auth.RoleAdmin)(admin.RemoveSemester),
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
