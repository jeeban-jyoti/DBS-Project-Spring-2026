package student

import (
	"database/sql"
	"encoding/json"
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

	// ✅ Added all columns from 'book' and 'category'
	// ✅ Added GROUP_CONCAT for subcategories and keywords
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
		var auths, subs, kws sql.NullString // Use NullString to prevent errors on empty joins

		err := rows.Scan(
			&b.BookID, &b.Title, &b.ISBN, &b.Publisher, &b.PublicationDate,
			&b.Edition, &b.Language, &b.Format, &b.Type, &b.PurchaseOption,
			&b.Price, &b.Quantity, &b.Category, &auths, &subs, &kws,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Convert pipe-separated strings from GROUP_CONCAT into Go slices
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

	// Use a temporary struct to handle the raw SQL GROUP_CONCAT strings
	var raw struct {
		BookFullResponse
		Authors       sql.NullString
		Subcategories sql.NullString
		Keywords      sql.NullString
		AvgRating     float64
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
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Final Response Object
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

	// Fetch Reviews
	rows, err := database.DB.Query(`SELECT student_id, rating, review_text, review_date FROM review WHERE book_id = ?`, bookID)
	if err == nil {
		defer rows.Close()
		var reviews []Review
		for rows.Next() {
			var rev Review
			if err := rows.Scan(&rev.StudentID, &rev.Rating, &rev.Text, &rev.Date); err == nil {
				reviews = append(reviews, rev)
			}
		}
		// You can add the Reviews slice to your BookFullResponse struct or a wrapper
		_ = reviews // Add to response wrapper if needed
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
