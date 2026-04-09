package authentication

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var u UserHTTPReq

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var user UserDBData

	err = database.DB.QueryRow(`
		SELECT 
			u.email,
			u.password_hash,
			CONCAT(u.first_name, ' ', u.last_name) as name,
			CASE 
				WHEN sa.employee_id IS NOT NULL THEN 'superadmin'
				WHEN a.employee_id IS NOT NULL THEN 'admin'
				WHEN cs.employee_id IS NOT NULL THEN 'support'
				WHEN s.student_id IS NOT NULL THEN 'student'
				ELSE 'unknown'
			END as role
		FROM user u
		LEFT JOIN student s ON u.user_id = s.student_id
		LEFT JOIN employee e ON u.user_id = e.employee_id
		LEFT JOIN customer_support cs ON e.employee_id = cs.employee_id
		LEFT JOIN administrator a ON e.employee_id = a.employee_id
		LEFT JOIN super_admin sa ON a.employee_id = sa.employee_id
		WHERE u.email = ?
	`, u.Email).Scan(&user.Email, &user.Password, &user.Name, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if u.Password != user.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	sessionID := GenerateSessionID()
	Create(sessionID, user.Email, user.Role)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "login successful",
		"role":    user.Role,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Write([]byte("Logged out successfully"))
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PasswordChangeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userEmail := r.Header.Get("user_email")
	if userEmail == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var currentPassword string
	err := database.DB.QueryRow(
		"SELECT password_hash FROM user WHERE email = ?",
		userEmail,
	).Scan(&currentPassword)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if req.OldPassword != currentPassword {
		http.Error(w, "Old password incorrect", http.StatusUnauthorized)
		return
	}

	_, err = database.DB.Exec(
		"UPDATE user SET password_hash = ? WHERE email = ?",
		req.NewPassword,
		userEmail,
	)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully. Please login again.",
	})
}
