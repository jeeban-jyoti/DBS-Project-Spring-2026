package student

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func FetchAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	subcategories := r.URL.Query().Get("subcategories")
	keywords := r.URL.Query().Get("keywords")

	query := `
		SELECT 
			b.book_id,
			b.title,
			b.isbn,
			b.publisher,
			b.publication_date,
			b.edition,
			b.language,
			b.format,
			b.type,
			b.purchase_option,
			b.price,
			b.quantity,
			c.name as category_name,
			GROUP_CONCAT(DISTINCT a.name SEPARATOR '|') as authors,
			GROUP_CONCAT(DISTINCT sc.name SEPARATOR '|') as subcategories,
			GROUP_CONCAT(DISTINCT k.keyword SEPARATOR '|') as keywords
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN book_author ba ON b.book_id = ba.book_id
		LEFT JOIN author a ON ba.author_id = a.author_id
		LEFT JOIN book_subcategory bs ON b.book_id = bs.book_id
		LEFT JOIN subcategory sc ON bs.subcategory_id = sc.subcategory_id
		LEFT JOIN book_keyword bk ON b.book_id = bk.book_id
		LEFT JOIN keyword k ON bk.keyword_id = k.keyword_id
	`

	args := []interface{}{}

	if category != "" {
		query += " AND c.name = ?"
		args = append(args, category)
	}

	if subcategories != "" {
		subs := strings.Split(subcategories, ",")
		placeholders := strings.Repeat("?,", len(subs))
		placeholders = placeholders[:len(placeholders)-1]
		query += " AND sc.name IN (" + placeholders + ")"
		for _, s := range subs {
			args = append(args, strings.TrimSpace(s))
		}
	}

	if keywords != "" {
		kws := strings.Split(keywords, ",")
		placeholders := strings.Repeat("?,", len(kws))
		placeholders = placeholders[:len(placeholders)-1]
		query += " AND k.keyword IN (" + placeholders + ")"
		for _, k := range kws {
			args = append(args, strings.TrimSpace(k))
		}
	}

	query += " GROUP BY b.book_id"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var books []BookFullResponse

	for rows.Next() {
		var b BookFullResponse
		var auths, subs, kws sql.NullString

		err := rows.Scan(
			&b.BookID, &b.Title, &b.ISBN, &b.Publisher, &b.PublicationDate,
			&b.Edition, &b.Language, &b.Format, &b.Type, &b.PurchaseOption,
			&b.Price, &b.Quantity, &b.Category, &auths, &subs, &kws,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if auths.Valid && auths.String != "" {
			b.Authors = strings.Split(auths.String, "|")
		}
		if subs.Valid && subs.String != "" {
			b.Subcategories = strings.Split(subs.String, "|")
		}
		if kws.Valid && kws.String != "" {
			b.Keywords = strings.Split(kws.String, "|")
		}

		books = append(books, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func FetchBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	bookID := r.URL.Query().Get("bookid")
	if bookID == "" {
		http.Error(w, "bookid required", http.StatusBadRequest)
		return
	}

	var raw struct {
		BookFullResponse
		Authors       sql.NullString
		Subcategories sql.NullString
		Keywords      sql.NullString
	}

	err := database.DB.QueryRow(`
		SELECT 
			b.book_id, b.title, b.isbn, b.publisher, b.publication_date, 
			b.edition, b.language, b.format, b.type, b.purchase_option, 
			b.price, b.quantity, c.name,
			GROUP_CONCAT(DISTINCT a.name SEPARATOR '|'),
			GROUP_CONCAT(DISTINCT sc.name SEPARATOR '|'),
			GROUP_CONCAT(DISTINCT k.keyword SEPARATOR '|'),
			IFNULL(AVG(r.rating), 0)
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN book_author ba ON b.book_id = ba.book_id
		LEFT JOIN author a ON ba.author_id = a.author_id
		LEFT JOIN book_subcategory bsc ON b.book_id = bsc.book_id
		LEFT JOIN subcategory sc ON bsc.subcategory_id = sc.subcategory_id
		LEFT JOIN book_keyword bk ON b.book_id = bk.book_id
		LEFT JOIN keyword k ON bk.keyword_id = k.keyword_id
		LEFT JOIN review r ON b.book_id = r.book_id
		WHERE b.book_id = ?
		GROUP BY b.book_id
	`, bookID).Scan(
		&raw.BookID, &raw.Title, &raw.ISBN, &raw.Publisher, &raw.PublicationDate,
		&raw.Edition, &raw.Language, &raw.Format, &raw.Type, &raw.PurchaseOption,
		&raw.Price, &raw.Quantity, &raw.Category,
		&raw.Authors, &raw.Subcategories, &raw.Keywords, &raw.AvgRating,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Process the result into the final response object
	response := raw.BookFullResponse
	if raw.Authors.Valid {
		response.Authors = strings.Split(raw.Authors.String, "|")
	}
	if raw.Subcategories.Valid {
		response.Subcategories = strings.Split(raw.Subcategories.String, "|")
	}
	if raw.Keywords.Valid {
		response.Keywords = strings.Split(raw.Keywords.String, "|")
	}

	rows, err := database.DB.Query(`
		SELECT student_id, rating, review_text, review_date 
		FROM review 
		WHERE book_id = ? 
		ORDER BY review_date DESC
	`, bookID)

	if err != nil {
		fmt.Printf("Error fetching reviews: %v\n", err)
	} else {
		defer rows.Close()
		response.Reviews = []Review{}
		for rows.Next() {
			var rev Review
			if err := rows.Scan(&rev.StudentID, &rev.Rating, &rev.Text, &rev.Date); err == nil {
				response.Reviews = append(response.Reviews, rev)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func FetchAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`SELECT name FROM category`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var categories []string

	for rows.Next() {
		var name string
		rows.Scan(&name)
		categories = append(categories, name)
	}

	json.NewEncoder(w).Encode(categories)
}

func FetchAllSubcategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "category is required", 400)
		return
	}

	rows, err := database.DB.Query(`
		SELECT sc.name
		FROM subcategory sc
		JOIN category c ON sc.category_id = c.category_id
		WHERE c.name = ?
	`, category)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var subs []string

	for rows.Next() {
		var name string
		rows.Scan(&name)
		subs = append(subs, name)
	}

	json.NewEncoder(w).Encode(subs)
}

func AddReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "User email missing in header", http.StatusBadRequest)
		return
	}

	var studentID int
	err := database.DB.QueryRow(`
		SELECT s.student_id FROM student s
		JOIN user u ON s.student_id = u.user_id
		WHERE u.email = ?
	`, email).Scan(&studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	var req struct {
		BookID     int    `json:"book_id"`
		Rating     int    `json:"rating"`
		ReviewText string `json:"review_text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BookID == 0 {
		http.Error(w, "book_id is required", http.StatusBadRequest)
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO review (student_id, book_id, rating, review_text)
		VALUES (?, ?, ?, ?)
	`, studentID, req.BookID, req.Rating, req.ReviewText)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "You have already reviewed this book", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Review added successfully"})
}
