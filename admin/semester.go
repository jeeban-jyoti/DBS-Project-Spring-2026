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

	// ✅ Basic validation
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

	// ✅ Validate course
	var exists int
	err = tx.QueryRow(`SELECT COUNT(*) FROM course WHERE course_id = ?`, req.CourseID).Scan(&exists)
	if err != nil || exists == 0 {
		http.Error(w, "Invalid course_id", http.StatusBadRequest)
		return
	}

	// ✅ Validate instructor
	err = tx.QueryRow(`SELECT COUNT(*) FROM instructor WHERE instructor_id = ?`, req.InstructorID).Scan(&exists)
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
