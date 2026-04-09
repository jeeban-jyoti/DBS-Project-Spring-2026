package student

import (
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func FetchAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT 
			b.book_id,
			b.title,
			GROUP_CONCAT(a.name SEPARATOR ', ') as authors,
			b.publication_date,
			b.price,
			b.quantity
		FROM book b
		LEFT JOIN book_author ba ON b.book_id = ba.book_id
		LEFT JOIN author a ON ba.author_id = a.author_id
		WHERE b.quantity > 0
		GROUP BY b.book_id
	`)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Authors, &b.PublicationDate, &b.Price, &b.Quantity)
		if err != nil {
			http.Error(w, "Scan error", 500)
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
