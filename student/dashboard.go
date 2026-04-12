package student

import (
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

	// ✅ Read query params
	category := r.URL.Query().Get("category")
	subcategories := r.URL.Query().Get("subcategories") // comma-separated
	keywords := r.URL.Query().Get("keywords")           // comma-separated

	query := `
		SELECT 
			b.book_id,
			b.title,
			GROUP_CONCAT(DISTINCT a.name SEPARATOR ', ') as authors,
			b.publication_date,
			b.price,
			b.quantity
		FROM book b
		LEFT JOIN book_author ba ON b.book_id = ba.book_id
		LEFT JOIN author a ON ba.author_id = a.author_id
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN book_subcategory bs ON b.book_id = bs.book_id
		LEFT JOIN subcategory sc ON bs.subcategory_id = sc.subcategory_id
		LEFT JOIN book_keyword bk ON b.book_id = bk.book_id
		LEFT JOIN keyword k ON bk.keyword_id = k.keyword_id
		WHERE b.quantity > 0
	`

	args := []interface{}{}

	// ✅ Category filter
	if category != "" {
		query += " AND c.name = ?"
		args = append(args, category)
	}

	// ✅ Subcategories filter
	if subcategories != "" {
		subs := strings.Split(subcategories, ",")
		placeholders := strings.Repeat("?,", len(subs))
		placeholders = placeholders[:len(placeholders)-1]

		query += " AND sc.name IN (" + placeholders + ")"
		for _, s := range subs {
			args = append(args, strings.TrimSpace(s))
		}
	}

	// ✅ Keywords filter
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

	var books []Book

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Authors, &b.PublicationDate, &b.Price, &b.Quantity)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		books = append(books, b)
	}

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

	var b BookDetail

	err := database.DB.QueryRow(`
		SELECT 
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

			GROUP_CONCAT(DISTINCT a.name SEPARATOR ', ') as authors,
			GROUP_CONCAT(DISTINCT sc.name SEPARATOR ', ') as subcategories,

			IFNULL(AVG(r.rating), 0) as avg_rating

		FROM book b
		LEFT JOIN book_author ba ON b.book_id = ba.book_id
		LEFT JOIN author a ON ba.author_id = a.author_id

		LEFT JOIN book_subcategory bsc ON b.book_id = bsc.book_id
		LEFT JOIN subcategory sc ON bsc.subcategory_id = sc.subcategory_id

		LEFT JOIN review r ON b.book_id = r.book_id

		WHERE b.book_id = ?
		GROUP BY b.book_id
	`, bookID).Scan(
		&b.Title,
		&b.ISBN,
		&b.Publisher,
		&b.PublicationDate,
		&b.Edition,
		&b.Language,
		&b.Format,
		&b.Type,
		&b.PurchaseOption,
		&b.Price,
		&b.Quantity,
		&b.Authors,
		&b.SubCategories,
		&b.AvgRating, // ⭐ NEW
	)

	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query(`
		SELECT student_id, rating, review_text, review_date
		FROM review
		WHERE book_id = ?
	`, bookID)

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			var rev Review
			err := rows.Scan(&rev.StudentID, &rev.Rating, &rev.Text, &rev.Date)
			if err == nil {
				b.Reviews = append(b.Reviews, rev)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
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
