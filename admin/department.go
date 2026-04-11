package admin

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	pass := make([]byte, length)
	for i := range pass {
		pass[i] = charset[rng.Intn(len(charset))]
	}
	return string(pass)
}

// Department
func AddDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddDepartmentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// ✅ Validate university exists
	var exists int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM university WHERE university_id = ?
	`, req.UniversityID).Scan(&exists)

	if err != nil || exists == 0 {
		http.Error(w, "Invalid university_id", http.StatusBadRequest)
		return
	}

	// ✅ Insert department
	res, err := database.DB.Exec(`
		INSERT INTO department (name, university_id)
		VALUES (?, ?)
	`, req.Name, req.UniversityID)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Department added successfully",
		"department_id": id,
	})
}

func RemoveDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Department ID required", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec(`
		DELETE FROM department WHERE department_id = ?
	`, id)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Department removed successfully"))
}

// Course
func AddCourse(w http.ResponseWriter, r *http.Request) {
	var req CourseReq
	json.NewDecoder(r.Body).Decode(&req)

	res, err := database.DB.Exec(`
		INSERT INTO course (name, university_id, year, semester)
		VALUES (?, ?, ?, ?)
	`, req.Name, req.UniversityID, req.Year, req.Semester)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"course_id": id})
}

func RemoveCourse(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	res, _ := database.DB.Exec(`DELETE FROM course WHERE course_id = ?`, id)
	rows, _ := res.RowsAffected()

	if rows == 0 {
		http.Error(w, "Course not found", 404)
		return
	}

	w.Write([]byte("Course deleted"))
}

// Instructor
func AddInstructor(w http.ResponseWriter, r *http.Request) {
	var req InstructorReq
	json.NewDecoder(r.Body).Decode(&req)

	password := generatePassword(10)

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO user (first_name, last_name, email, address, phone, password_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.FirstName, req.LastName, req.Email, req.Address, req.Phone, password)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	userID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		INSERT INTO instructor (instructor_id, university_id, department_id)
		VALUES (?, ?, ?)
	`, userID, req.UniversityID, req.DepartmentID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Instructor added",
		"password": password,
	})
}

func RemoveInstructor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	tx.Exec(`DELETE FROM instructor WHERE instructor_id = ?`, id)
	tx.Exec(`DELETE FROM user WHERE user_id = ?`, id)

	tx.Commit()

	w.Write([]byte("Instructor removed"))
}

func UpdateInstructor(w http.ResponseWriter, r *http.Request) {
	var req InstructorReq
	id := r.URL.Query().Get("id")

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	tx.Exec(`
		UPDATE user 
		SET first_name=?, last_name=?, email=?, address=?, phone=?
		WHERE user_id=?
	`, req.FirstName, req.LastName, req.Email, req.Address, req.Phone, id)

	tx.Exec(`
		UPDATE instructor 
		SET university_id=?, department_id=?
		WHERE instructor_id=?
	`, req.UniversityID, req.DepartmentID, id)

	tx.Commit()

	w.Write([]byte("Instructor updated"))
}

// Student
func AddStudent(w http.ResponseWriter, r *http.Request) {
	var req StudentReq
	json.NewDecoder(r.Body).Decode(&req)

	password := generatePassword(10)

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	// user
	res, err := tx.Exec(`
		INSERT INTO user (first_name, last_name, email, address, phone, password_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.FirstName, req.LastName, req.Email, req.Address, req.Phone, password)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	userID, _ := res.LastInsertId()

	// student
	_, err = tx.Exec(`
		INSERT INTO student (student_id, date_of_birth, university_id, major, status, year_of_study)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, req.DOB, req.UniversityID, req.Major, req.Status, req.YearOfStudy)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Student added",
		"password": password,
	})
}

func RemoveStudent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	tx.Exec(`DELETE FROM student WHERE student_id = ?`, id)
	tx.Exec(`DELETE FROM user WHERE user_id = ?`, id)

	tx.Commit()

	w.Write([]byte("Student removed"))
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	var req StudentReq
	id := r.URL.Query().Get("id")

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	tx.Exec(`
		UPDATE user 
		SET first_name=?, last_name=?, email=?, address=?, phone=?
		WHERE user_id=?
	`, req.FirstName, req.LastName, req.Email, req.Address, req.Phone, id)

	tx.Exec(`
		UPDATE student 
		SET date_of_birth=?, university_id=?, major=?, status=?, year_of_study=?
		WHERE student_id=?
	`, req.DOB, req.UniversityID, req.Major, req.Status, req.YearOfStudy, id)

	tx.Commit()

	w.Write([]byte("Student updated"))
}
