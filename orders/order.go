package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

// cartItem is a small helper struct to hold data temporarily
// to avoid "busy buffer" errors during transactions.

func getStudentID(email string) (int, error) {
	var id int
	err := database.DB.QueryRow(`
		SELECT s.student_id FROM student s
		JOIN user u ON s.student_id = u.user_id
		WHERE u.email = ?
	`, email).Scan(&id)
	return id, err
}

// Helper function to handle both Buy and Borrow logic
func processOrder(w http.ResponseWriter, r *http.Request, purchaseType string) {
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "User email missing in header", http.StatusBadRequest)
		return
	}

	studentID, err := getStudentID(email)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Start Transaction
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Defer a rollback. If tx.Commit() is called, the rollback does nothing.
	// If the function returns early due to an error, this ensures the DB stays consistent.
	defer tx.Rollback()

	// 1. Create the Order
	res, err := tx.Exec("INSERT INTO `order`(student_id, status) VALUES(?, 'new')", studentID)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	orderID, _ := res.LastInsertId()

	// 2. Get items from Cart and STORE THEM IN A SLICE
	// We do this so we can close the 'rows' cursor before performing updates.
	var itemsToProcess []cartItem
	rows, err := tx.Query(`
		SELECT ci.book_id, ci.quantity
		FROM cart c
		JOIN cart_item ci ON c.cart_id = ci.cart_id
		WHERE c.student_id = ?
	`, studentID)
	if err != nil {
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var item cartItem
		if err := rows.Scan(&item.bookID, &item.qty); err != nil {
			rows.Close()
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}
		itemsToProcess = append(itemsToProcess, item)
	}
	rows.Close() // <--- CRITICAL: Close the cursor NOW. Connection is now free for Exec.

	if len(itemsToProcess) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// 3. Process each item from the slice
	for _, item := range itemsToProcess {
		// Insert into order_item
		_, err = tx.Exec(`
			INSERT INTO order_item(order_id, book_id, quantity, purchase_type)
			VALUES(?, ?, ?, ?)
		`, orderID, item.bookID, item.qty, purchaseType)
		if err != nil {
			http.Error(w, "Failed to add item to order", http.StatusInternalServerError)
			return
		}

		// DECREASE QUANTITY (Required for both Buy and Rent)
		// The WHERE clause prevents quantity from going below 0 (Atomic check)
		updateRes, err := tx.Exec(`
			UPDATE book 
			SET quantity = quantity - ? 
			WHERE book_id = ? AND quantity >= ?
		`, item.qty, item.bookID, item.qty)
		if err != nil {
			http.Error(w, "Database error updating stock", http.StatusInternalServerError)
			return
		}

		// Check if the book was actually available (if 0 rows affected, it means quantity < item.qty)
		rowsAffected, _ := updateRes.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, fmt.Sprintf("Book ID %d is out of stock or insufficient quantity", item.bookID), http.StatusBadRequest)
			return // This triggers the deferred Rollback()
		}
	}

	// 4. Clear the Cart
	_, err = tx.Exec(`
		DELETE ci FROM cart_item ci
		JOIN cart c ON ci.cart_id = c.cart_id
		WHERE c.student_id = ?
	`, studentID)
	if err != nil {
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to finalize order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%s order placed successfully", purchaseType)})
}

func PlaceBuyOrder(w http.ResponseWriter, r *http.Request) {
	processOrder(w, r, "buy")
}

func PlaceBorrowOrder(w http.ResponseWriter, r *http.Request) {
	processOrder(w, r, "rent")
}

// --- Response Structs ---

// Helper to fetch book details for a specific order
func getOrderItems(tx *sql.Tx, orderID int) ([]OrderItemDetail, error) {
	rows, err := tx.Query(`
		SELECT b.title, oi.quantity, oi.purchase_type 
		FROM order_item oi 
		JOIN book b ON oi.book_id = b.book_id 
		WHERE oi.order_id = ?
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItemDetail
	for rows.Next() {
		var item OrderItemDetail
		if err := rows.Scan(&item.BookTitle, &item.Quantity, &item.PurchaseType); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// --- API Endpoints ---

// ShowOrder: GET /api/v1/order?id=123
func ShowOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid Order ID", http.StatusBadRequest)
		return
	}

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	var o OrderDetail
	err = tx.QueryRow(`
		SELECT order_id, student_id, created_date, fulfilled_date, shipping_type, card_holder_name, status 
		FROM `+"`order`"+` WHERE order_id = ?
	`, orderID).Scan(&o.OrderID, &o.StudentID, &o.CreatedDate, &o.FulfilledDate, &o.ShippingType, &o.CardHolderName, &o.Status)

	if err == sql.ErrNoRows {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	items, err := getOrderItems(tx, orderID)
	if err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}
	o.Items = items

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// ShowOrders: GET /api/v1/orders
func ShowOrders(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	studentID, err := getStudentID(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. First, fetch only the order headers and store them in a slice
	rows, err := tx.Query(`
		SELECT order_id, student_id, created_date, fulfilled_date, shipping_type, card_holder_name, status 
		FROM `+"`order`"+` 
		WHERE student_id = ? AND status NOT IN ('shipped', 'canceled')
	`, studentID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var orders []OrderDetail
	for rows.Next() {
		var o OrderDetail
		err := rows.Scan(&o.OrderID, &o.StudentID, &o.CreatedDate, &o.FulfilledDate, &o.ShippingType, &o.CardHolderName, &o.Status)
		if err != nil {
			rows.Close()
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}
		o.StudentEmail = email // Set the email from the header
		orders = append(orders, o)
	}

	// 2. CRITICAL: Close the rows cursor NOW.
	// This tells MySQL "I am done reading the list of orders,"
	// and clears the buffer for the next set of queries.
	rows.Close()

	// 3. Now that the connection is free, fetch the items for each order
	for i := range orders {
		items, err := getOrderItems(tx, orders[i].OrderID)
		if err != nil {
			// We log the error but continue so one bad order doesn't crash the whole list
			fmt.Printf("Error fetching items for order %d: %v\n", orders[i].OrderID, err)
			continue
		}
		orders[i].Items = items
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to finalize transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func ShowAllOrders(w http.ResponseWriter, r *http.Request) {
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. First, fetch ALL orders and store them in a slice
	rows, err := tx.Query(`
		SELECT o.order_id, o.student_id, u.email, o.created_date, o.fulfilled_date, o.shipping_type, o.card_holder_name, o.status 
		FROM ` + "`order`" + ` o
		JOIN user u ON o.student_id = u.user_id
		ORDER BY o.created_date DESC
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var orders []OrderDetail
	for rows.Next() {
		var o OrderDetail
		err := rows.Scan(
			&o.OrderID,
			&o.StudentID,
			&o.StudentEmail,
			&o.CreatedDate,
			&o.FulfilledDate,
			&o.ShippingType,
			&o.CardHolderName,
			&o.Status,
		)
		if err != nil {
			rows.Close()
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	// 2. CRITICAL: Close the rows cursor NOW.
	// The connection is now "IDLE" and ready for new queries.
	rows.Close()

	// 3. Now that the connection is free, loop through the slice
	// and fetch the items for each order.
	for i := range orders {
		items, err := getOrderItems(tx, orders[i].OrderID)
		if err != nil {
			// We log the error but continue processing other orders
			fmt.Printf("Error fetching items for order %d: %v\n", orders[i].OrderID, err)
			continue
		}
		orders[i].Items = items
	}

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GenerateOrderCancellation: POST /api/v1/cancelOrder
func GenerateOrderCancellation(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	studentID, err := getStudentID(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	orderIDStr := r.URL.Query().Get("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid Order ID", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Check if order exists and if it's already shipped/cancelled
	var status string
	err = tx.QueryRow(`SELECT status FROM `+"`order`"+` WHERE order_id = ? AND student_id = ?`, orderID, studentID).Scan(&status)
	if err == sql.ErrNoRows {
		http.Error(w, "Order not found or not owned by user", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if status == "shipped" || status == "canceled" {
		http.Error(w, "Cannot cancel an order that is already shipped or cancelled", http.StatusBadRequest)
		return
	}

	// 2. COLLECT the items to be returned
	type itemReturn struct {
		bookID int
		qty    int
	}
	var itemsToReturn []itemReturn

	rows, err := tx.Query(`SELECT book_id, quantity FROM order_item WHERE order_id = ?`, orderID)
	if err != nil {
		http.Error(w, "Failed to fetch order items", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var ir itemReturn
		if err := rows.Scan(&ir.bookID, &ir.qty); err != nil {
			rows.Close()
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}
		itemsToReturn = append(itemsToReturn, ir)
	}

	// CRITICAL: Close the rows cursor NOW so the connection is free for UPDATES
	rows.Close()

	// 3. PROCESS the inventory return
	for _, item := range itemsToReturn {
		_, err = tx.Exec(`UPDATE book SET quantity = quantity + ? WHERE book_id = ?`, item.qty, item.bookID)
		if err != nil {
			http.Error(w, "Failed to return stock for book ID "+fmt.Sprint(item.bookID), http.StatusInternalServerError)
			return
		}
	}

	// 4. Update status to cancelled
	_, err = tx.Exec(`UPDATE `+"`order`"+` SET status = 'canceled' WHERE order_id = ?`, orderID)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to finalize cancellation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order cancelled and stock returned"})
}

// ChangeOrderStatus: PUT /api/v1/updateOrderStatus
func ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Typically an Admin API
	orderIDStr := r.URL.Query().Get("id")
	newStatus := r.URL.Query().Get("status")
	orderID, _ := strconv.Atoi(orderIDStr)

	// Validate status input
	validStatuses := map[string]bool{
		"new": true, "processed": true, "awaiting_shipping": true, "shipped": true, "canceled": true,
	}
	if !validStatuses[newStatus] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`UPDATE `+"`order`"+` SET status = ? WHERE order_id = ?`, newStatus, orderID)
	if err != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated successfully"})
}
