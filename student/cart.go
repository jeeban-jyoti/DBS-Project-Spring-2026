package student

import (
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func getStudentID(email string) (int, error) {
	var id int
	err := database.DB.QueryRow(`
		SELECT student_id FROM student s
		JOIN user u ON s.student_id = u.user_id
		WHERE u.email = ?
	`, email).Scan(&id)
	return id, err
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	var req AddToCartReq
	json.NewDecoder(r.Body).Decode(&req)

	email := r.Header.Get("user_email")
	studentID, _ := getStudentID(email)

	var cartID int
	err := database.DB.QueryRow("SELECT cart_id FROM cart WHERE student_id = ?", studentID).Scan(&cartID)

	if err != nil {
		res, _ := database.DB.Exec("INSERT INTO cart(student_id) VALUES(?)", studentID)
		id, _ := res.LastInsertId()
		cartID = int(id)
	}

	_, err = database.DB.Exec(`
		INSERT INTO cart_item(cart_id, book_id, quantity)
		VALUES(?, ?, ?)
		ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)
	`, cartID, req.BookID, req.Quantity)

	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "added to cart"})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var req RemoveFromCartReq
	json.NewDecoder(r.Body).Decode(&req)

	email := r.Header.Get("user_email")
	studentID, _ := getStudentID(email)

	_, err := database.DB.Exec(`
		DELETE ci FROM cart_item ci
		JOIN cart c ON ci.cart_id = c.cart_id
		WHERE c.student_id = ? AND ci.book_id = ?
	`, studentID, req.BookID)

	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "removed"})
}

func ShowCart(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	studentID, _ := getStudentID(email)

	rows, _ := database.DB.Query(`
		SELECT b.book_id, b.title, b.price, ci.quantity
		FROM cart c
		JOIN cart_item ci ON c.cart_id = ci.cart_id
		JOIN book b ON ci.book_id = b.book_id
		WHERE c.student_id = ?
	`, studentID)

	var books []Book

	for rows.Next() {
		var b Book
		rows.Scan(&b.ID, &b.Title, &b.Price, &b.Quantity)
		books = append(books, b)
	}

	json.NewEncoder(w).Encode(books)
}

func PlaceBuyOrder(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	studentID, _ := getStudentID(email)

	tx, _ := database.DB.Begin()

	res, _ := tx.Exec("INSERT INTO `order`(student_id, status) VALUES(?, 'new')", studentID)
	orderID, _ := res.LastInsertId()

	rows, _ := tx.Query(`
		SELECT ci.book_id, ci.quantity
		FROM cart c
		JOIN cart_item ci ON c.cart_id = ci.cart_id
		WHERE c.student_id = ?
	`, studentID)

	for rows.Next() {
		var bookID, qty int
		rows.Scan(&bookID, &qty)

		tx.Exec(`
			INSERT INTO order_item(order_id, book_id, quantity, purchase_type)
			VALUES(?, ?, ?, 'buy')
		`, orderID, bookID, qty)

		tx.Exec(`
			UPDATE book SET quantity = quantity - ?
			WHERE book_id = ? AND quantity >= ?
		`, qty, bookID, qty)
	}

	tx.Exec(`
		DELETE ci FROM cart_item ci
		JOIN cart c ON ci.cart_id = c.cart_id
		WHERE c.student_id = ?
	`, studentID)

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{"message": "order placed"})
}

func PlaceBorrowOrder(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	studentID, _ := getStudentID(email)

	tx, _ := database.DB.Begin()

	res, _ := tx.Exec("INSERT INTO `order`(student_id, status) VALUES(?, 'new')", studentID)
	orderID, _ := res.LastInsertId()

	rows, _ := tx.Query(`
		SELECT ci.book_id, ci.quantity
		FROM cart c
		JOIN cart_item ci ON c.cart_id = ci.cart_id
		WHERE c.student_id = ?
	`, studentID)

	for rows.Next() {
		var bookID, qty int
		rows.Scan(&bookID, &qty)

		tx.Exec(`
			INSERT INTO order_item(order_id, book_id, quantity, purchase_type)
			VALUES(?, ?, ?, 'rent')
		`, orderID, bookID, qty)
	}

	tx.Exec(`
		DELETE ci FROM cart_item ci
		JOIN cart c ON ci.cart_id = c.cart_id
		WHERE c.student_id = ?
	`, studentID)

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{"message": "borrow placed"})
}
