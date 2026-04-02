package authentication

import (
	"net/http"
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

		// Attach to request
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
