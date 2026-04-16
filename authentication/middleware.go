package authentication

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

const (
	RoleStudent    = "student"
	RoleSupport    = "support"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "superadmin"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, ok := Get(cookie.Value)
		if !ok {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_email", session.Email)
		r.Header.Set("user_role", session.Role)

		next(w, r)
	}
}

func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			role := r.Header.Get("user_role")

			for _, allowed := range allowedRoles {
				if role == allowed {
					next(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}

func LogRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("user_email")
		role := r.Header.Get("user_role")

		// fmt.Printf("Logging Request: User=%s, Role=%s, URL=%s\n", email, role, r.URL.Path)

		if email == "" || role == "" {
			fmt.Println("Logging skipped: User email or role is empty")
			next(w, r)
			return
		}

		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString := string(bodyBytes)

		tx, err := database.DB.Begin()
		if err != nil {
			fmt.Println("DB Transaction Error:", err)
			next(w, r)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`INSERT INTO request_logs (user_email, user_role, url, body) VALUES (?, ?, ?, ?)`,
			email, role, r.URL.String(), bodyString)

		if err != nil {
			fmt.Println("SQL Insert Error:", err)
			next(w, r)
			return
		}

		if err := tx.Commit(); err != nil {
			fmt.Println("Commit Error:", err)
		}

		next(w, r)
	}
}
