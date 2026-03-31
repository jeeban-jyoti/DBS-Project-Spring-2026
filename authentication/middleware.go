package authentication

import (
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email, ok := Get(cookie.Value)
		if !ok {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_email", email)

		next(w, r)
	}
}
