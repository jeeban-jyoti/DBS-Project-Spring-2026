package admin

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	pass := make([]byte, length)
	for i := range pass {
		pass[i] = charset[rng.Intn(len(charset))]
	}
	return string(pass)
}

func AddStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddStudentsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Students) == 0 {
		http.Error(w, "No students provided", http.StatusBadRequest)
		return
	}

	password := generatePassword(10)

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	for _, s := range req.Students {
		_, err := tx.Exec(
			"INSERT INTO users (email, passwordHash, type, name) VALUES (?, ?, 'student', ?)",
			s.Email, password, s.Name,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Students added successfully",
		"count":    len(req.Students),
		"password": password,
	})
}

func RemoveStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RemoveStudentsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Emails) == 0 {
		http.Error(w, "No emails provided", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	for _, email := range req.Emails {
		_, err := tx.Exec(
			"DELETE FROM users WHERE email = ? AND type = 'student'",
			email,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Students removed successfully",
		"count":   len(req.Emails),
	})
}

func AddBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddBooksReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Books) == 0 {
		http.Error(w, "No books provided", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	for _, b := range req.Books {
		_, err := tx.Exec(
			`INSERT INTO books (id, title, author, pub_year, status, quantity)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			b.ID, b.Title, b.Author, b.PubYear, b.Status, b.Quantity,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Books added successfully",
		"count":   len(req.Books),
	})
}

func RemoveBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RemoveBooksReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "No book IDs provided", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	for _, id := range req.IDs {
		_, err := tx.Exec(
			"DELETE FROM books WHERE id = ?",
			id,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Books removed successfully",
		"count":   len(req.IDs),
	})
}
