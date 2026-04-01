package student

import (
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var request AddToCartReq

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO carts(studentId, bookId, quantity) VALUES(?, ?, ?)", request.UserID, request.BookID, request.Quantity)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "added to cart",
	})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var request RemoveFromCartReq

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		"DELETE FROM carts WHERE studentId = ? AND bookId = ?",
		request.UserID,
		request.BookID,
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Item not found in cart", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "removed from cart",
	})
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userid")
	if userID == "" {
		http.Error(w, "userid required", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(
		`SELECT b.id, b.title, b.author, b.pub_year, b.status, c.quantity
		 FROM carts c
		 JOIN books b ON c.bookId = b.id
		 WHERE c.studentId = ?`,
		userID,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.PubYear,
			&b.Status,
			&b.Quantity,
		)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}

	json.NewEncoder(w).Encode(books)
}

func PlaceBuyOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OrderPlaceReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", 500)
		return
	}

	rows, err := tx.Query("SELECT bookId, quantity FROM carts WHERE studentId = ?", req.UserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bookID, qty int
		rows.Scan(&bookID, &qty)

		result, err := tx.Exec(
			"UPDATE books SET quantity = quantity - ? WHERE id = ? AND quantity >= ?",
			qty, bookID, qty,
		)

		rowsAffected, _ := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			tx.Rollback()
			http.Error(w, "Out of stock", 400)
			return
		}
	}

	_, err = tx.Exec("DELETE FROM carts WHERE studentId = ?", req.UserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error clearing cart", 500)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message": "order placed successfully",
	})
}

func PlaceBorrowOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OrderPlaceReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", 500)
		return
	}

	rows, err := tx.Query("SELECT bookId FROM carts WHERE studentId = ?", req.UserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bookID int
		rows.Scan(&bookID)

		_, err := tx.Exec(
			"UPDATE books SET status = 'issued' WHERE id = ?",
			bookID,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error issuing book", 500)
			return
		}
	}

	_, err = tx.Exec("DELETE FROM carts WHERE studentId = ?", req.UserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error clearing cart", 500)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message": "borrow successful",
	})
}
