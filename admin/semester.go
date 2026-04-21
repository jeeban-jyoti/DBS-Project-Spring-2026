package admin

import (
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func AddSemester(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddSemesterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Year == 0 || req.Season == "" || req.CourseID == 0 || req.InstructorID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var exists int
	err = tx.QueryRow(`SELECT COUNT(*) FROM university WHERE university_id = ?`, req.UniversityID).Scan(&exists)
	if err != nil || exists == 0 {
		http.Error(w, "Invalid university_id", http.StatusBadRequest)
		return
	}

	err = tx.QueryRow(`SELECT COUNT(*) FROM course WHERE course_id = ? AND university_id = ?`, req.CourseID, req.UniversityID).Scan(&exists)
	if err != nil || exists == 0 {
		http.Error(w, "Invalid course_id", http.StatusBadRequest)
		return
	}

	err = tx.QueryRow(`SELECT COUNT(*) FROM instructor WHERE instructor_id = ? AND university_id = ?`, req.InstructorID, req.UniversityID).Scan(&exists)
	if err != nil || exists == 0 {
		http.Error(w, "Invalid instructor_id", http.StatusBadRequest)
		return
	}

	// 1. Insert semester
	res, err := tx.Exec(`
		INSERT INTO semester (year, season, course_id, instructor_id, university_id)
		VALUES (?, ?, ?, ?, ?)
	`, req.Year, req.Season, req.CourseID, req.InstructorID, req.UniversityID)

	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	semID64, _ := res.LastInsertId()
	semID := int(semID64)

	// 2. Link books (optional)
	for _, bookID := range req.BookIDs {
		// validate book
		err := tx.QueryRow(`SELECT COUNT(*) FROM book WHERE book_id = ?`, bookID).Scan(&exists)
		if err != nil || exists == 0 {
			http.Error(w, "Invalid book_id: "+string(rune(bookID)), http.StatusBadRequest)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO semester_book (sem_id, book_id)
			VALUES (?, ?)
		`, semID, bookID)

		if err != nil {
			http.Error(w, "Failed to link book", http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Semester added successfully",
		"sem_id":  semID,
	})
}

func RemoveSemester(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "sem_id required", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Delete mappings
	tx.Exec(`DELETE FROM semester_book WHERE sem_id = ?`, id)

	// 2. Delete semester
	res, err := tx.Exec(`DELETE FROM semester WHERE sem_id = ?`, id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Semester not found", http.StatusNotFound)
		return
	}

	tx.Commit()

	w.Write([]byte("Semester deleted successfully"))
}

func FetchSemesters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT 
			s.sem_id,
			s.year,
			s.season,
			c.name,
			CONCAT(i.first_name, ' ', i.last_name) as instructor,
			u.name as university
		FROM semester s
		JOIN course c ON s.course_id = c.course_id
		JOIN instructor i ON s.instructor_id = i.instructor_id
		JOIN university u ON s.university_id = u.university_id
	`)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var semesters []SemesterDetail

	for rows.Next() {
		var sem SemesterDetail

		err := rows.Scan(
			&sem.SemID,
			&sem.Year,
			&sem.Season,
			&sem.CourseName,
			&sem.InstructorName,
			&sem.UniversityName,
		)
		if err != nil {
			http.Error(w, "Scan error", 500)
			return
		}

		// 🔥 Fetch books for each semester
		bookRows, _ := database.DB.Query(`
			SELECT b.title
			FROM semester_book sb
			JOIN book b ON sb.book_id = b.book_id
			WHERE sb.sem_id = ?
		`, sem.SemID)

		var books []string
		for bookRows.Next() {
			var title string
			bookRows.Scan(&title)
			books = append(books, title)
		}
		bookRows.Close()

		sem.Books = books

		semesters = append(semesters, sem)
	}

	json.NewEncoder(w).Encode(semesters)
}
