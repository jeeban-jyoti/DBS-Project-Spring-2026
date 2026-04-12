package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	var req AddBookReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. CATEGORY
	var categoryID int
	err = tx.QueryRow(`SELECT category_id FROM category WHERE name = ?`, req.Category).Scan(&categoryID)
	if err == sql.ErrNoRows {
		res, _ := tx.Exec(`INSERT INTO category (name) VALUES (?)`, req.Category)
		id, _ := res.LastInsertId()
		categoryID = int(id)
	}

	// 2. INSERT BOOK
	res, err := tx.Exec("INSERT INTO book (title, isbn, publisher, publication_date, edition, language, format, `type`, purchase_option, price, quantity, category_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Title, req.ISBN, req.Publisher, req.PublicationDate,
		req.Edition, req.Language, req.Format, req.Type,
		req.PurchaseOption, req.Price, req.Quantity, categoryID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookID64, _ := res.LastInsertId()
	bookID := int(bookID64)

	// 3. AUTHORS
	for _, name := range req.Authors {
		var authorID int
		err := tx.QueryRow(`SELECT author_id FROM author WHERE name = ?`, name).Scan(&authorID)
		if err == sql.ErrNoRows {
			res, _ := tx.Exec(`INSERT INTO author (name) VALUES (?)`, name)
			id, _ := res.LastInsertId()
			authorID = int(id)
		}
		tx.Exec(`INSERT INTO book_author (book_id, author_id) VALUES (?, ?)`, bookID, authorID)
	}

	// 4. KEYWORDS
	for _, kw := range req.Keywords {
		var keywordID int
		err := tx.QueryRow(`SELECT keyword_id FROM keyword WHERE keyword = ?`, kw).Scan(&keywordID)
		if err == sql.ErrNoRows {
			res, _ := tx.Exec(`INSERT INTO keyword (keyword) VALUES (?)`, kw)
			id, _ := res.LastInsertId()
			keywordID = int(id)
		}
		tx.Exec(`INSERT INTO book_keyword (book_id, keyword_id) VALUES (?, ?)`, bookID, keywordID)
	}

	// 5. SUBCATEGORIES
	for _, sub := range req.Subcategories {
		var subID int
		err := tx.QueryRow(`SELECT subcategory_id FROM subcategory WHERE name = ?`, sub).Scan(&subID)
		if err == sql.ErrNoRows {
			res, _ := tx.Exec(`INSERT INTO subcategory (name, category_id) VALUES (?, ?)`, sub, categoryID)
			id, _ := res.LastInsertId()
			subID = int(id)
		}
		tx.Exec(`INSERT INTO book_subcategory (book_id, subcategory_id) VALUES (?, ?)`, bookID, subID)
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Book added successfully",
	})
}

func RemoveBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Book ID required", http.StatusBadRequest)
		return
	}

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	// Delete mappings first
	tx.Exec(`DELETE FROM book_author WHERE book_id = ?`, id)
	tx.Exec(`DELETE FROM book_keyword WHERE book_id = ?`, id)
	tx.Exec(`DELETE FROM book_subcategory WHERE book_id = ?`, id)

	// Delete book
	res, _ := tx.Exec(`DELETE FROM book WHERE book_id = ?`, id)

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// CLEANUP: AUTHORS
	tx.Exec(`
		DELETE FROM author
		WHERE author_id NOT IN (SELECT DISTINCT author_id FROM book_author)
	`)

	// CLEANUP: KEYWORDS
	tx.Exec(`
		DELETE FROM keyword
		WHERE keyword_id NOT IN (SELECT DISTINCT keyword_id FROM book_keyword)
	`)

	// CLEANUP: SUBCATEGORIES
	tx.Exec(`
		DELETE FROM subcategory
		WHERE subcategory_id NOT IN (SELECT DISTINCT subcategory_id FROM book_subcategory)
	`)

	// CLEANUP: CATEGORY
	tx.Exec(`
		DELETE FROM category
		WHERE category_id NOT IN (SELECT DISTINCT category_id FROM book)
	`)

	tx.Commit()

	w.Write([]byte("Book deleted successfully"))
}
