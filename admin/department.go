package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Delete mappings
	_, err = tx.Exec(`
		DELETE FROM course_department WHERE department_id = ?
	`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Delete orphan courses
	_, err = tx.Exec(`
		DELETE FROM course
		WHERE course_id NOT IN (
			SELECT DISTINCT course_id FROM course_department
		)
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Delete department
	res, err := tx.Exec(`
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

	// Commit
	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Department removed successfully (orphan courses cleaned)"))
}

func GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	univID := r.URL.Query().Get("university_id")

	var rows *sql.Rows
	var err error

	if univID != "" {
		rows, err = database.DB.Query(`
			SELECT d.department_id, d.name, u.name
			FROM department d
			JOIN university u ON d.university_id = u.university_id
			WHERE d.university_id = ?
		`, univID)
	} else {
		rows, err = database.DB.Query(`
			SELECT d.department_id, d.name, u.name
			FROM department d
			JOIN university u ON d.university_id = u.university_id
		`)
	}

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Department struct {
		ID           int    `json:"department_id"`
		Name         string `json:"name"`
		UniversityName string    `json:"university_name"`
	}

	var departments []Department

	for rows.Next() {
		var d Department
		if err := rows.Scan(&d.ID, &d.Name, &d.UniversityName); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		departments = append(departments, d)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}

// Course
func AddCourse(w http.ResponseWriter, r *http.Request) {
	var req CourseReq

	// Decode request
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

	// ✅ Validate university exists
	var exists int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM university WHERE university_id = ?
	`, req.UniversityID).Scan(&exists)

	if err != nil || exists == 0 {
		http.Error(w, "Invalid university_id", http.StatusBadRequest)
		return
	}

	// 1. Insert course
	res, err := tx.Exec(`
		INSERT INTO course (name, university_id, year)
		VALUES (?, ?, ?)
	`, req.Name, req.UniversityID, req.Year)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	courseID64, _ := res.LastInsertId()
	courseID := int(courseID64)

	// 2. Insert course_department mappings
	for _, deptID := range req.Departments {

		// ✅ Validate department belongs to same university
		err := tx.QueryRow(`
			SELECT COUNT(*) FROM department 
			WHERE department_id = ? AND university_id = ?
		`, deptID, req.UniversityID).Scan(&exists)

		if err != nil || exists == 0 {
			http.Error(w, "Invalid department_id for this university", http.StatusBadRequest)
			return
		}

		// Insert mapping
		_, err = tx.Exec(`
			INSERT INTO course_department (course_id, department_id)
			VALUES (?, ?)
		`, courseID, deptID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Course added successfully",
		"course_id": courseID,
	})
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// Optional: validate university exists
	var exists int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM university WHERE university_id = ?
	`, req.UniversityID).Scan(&exists)

	if err != nil || exists == 0 {
		http.Error(w, "Invalid university_id", 400)
		return
	}

	// Optional: validate department belongs to university
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM department 
		WHERE department_id = ? AND university_id = ?
	`, req.DepartmentID, req.UniversityID).Scan(&exists)

	if err != nil || exists == 0 {
		http.Error(w, "Invalid department_id", 400)
		return
	}

	// ✅ Insert directly into instructor
	res, err := database.DB.Exec(`
		INSERT INTO instructor (first_name, last_name, university_id, department_id)
		VALUES (?, ?, ?, ?)
	`, req.FirstName, req.LastName, req.UniversityID, req.DepartmentID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := res.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Instructor added",
		"instructor_id": id,
	})
}

func RemoveInstructor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Instructor ID required", 400)
		return
	}

	res, err := database.DB.Exec(`
		DELETE FROM instructor WHERE instructor_id = ?
	`, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Instructor not found", 404)
		return
	}

	w.Write([]byte("Instructor removed"))
}

func UpdateInstructor(w http.ResponseWriter, r *http.Request) {
	var req InstructorReq
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Instructor ID required", 400)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// ✅ Validate university exists
	var exists int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM university WHERE university_id = ?
	`, req.UniversityID).Scan(&exists)

	if err != nil || exists == 0 {
		http.Error(w, "Invalid university_id", 400)
		return
	}

	// ✅ Validate department belongs to university
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM department 
		WHERE department_id = ? AND university_id = ?
	`, req.DepartmentID, req.UniversityID).Scan(&exists)

	if err != nil || exists == 0 {
		http.Error(w, "Invalid department_id", 400)
		return
	}

	_, err = database.DB.Exec(`
		UPDATE instructor 
		SET first_name=?, last_name=?, university_id=?, department_id=?
		WHERE instructor_id=?
	`, req.FirstName, req.LastName, req.UniversityID, req.DepartmentID, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("Instructor updated"))
}

func FetchInstructors(w http.ResponseWriter, r *http.Request) {
	// 1. Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Execute Query with JOINs to get human-readable names
	rows, err := database.DB.Query(`
		SELECT 
			i.instructor_id, 
			i.first_name, 
			i.last_name, 
			u.name as university_name, 
			d.name as department_name
		FROM instructor i
		JOIN university u ON i.university_id = u.university_id
		JOIN department d ON i.department_id = d.department_id
		ORDER BY i.last_name ASC
	`)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 3. Iterate through the results
	var instructors []InstructorDetail

	for rows.Next() {
		var inst InstructorDetail
		err := rows.Scan(
			&inst.InstructorID,
			&inst.FirstName,
			&inst.LastName,
			&inst.UniversityName,
			&inst.DepartmentName,
		)
		if err != nil {
			http.Error(w, "Data scan error", http.StatusInternalServerError)
			return
		}
		instructors = append(instructors, inst)
	}

	// 4. Handle empty state (return empty list instead of null)
	if instructors == nil {
		instructors = []InstructorDetail{}
	}

	// 5. Return as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instructors)
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

	fmt.Println(req.DOB)

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
	json.NewDecoder(r.Body).Decode(&req)
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

func FetchCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT 
			c.course_id,
			c.name,
			c.year,
			u.name
		FROM course c
		JOIN university u ON c.university_id = u.university_id
	`)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var courses []CourseDetail

	for rows.Next() {
		var c CourseDetail

		err := rows.Scan(
			&c.CourseID,
			&c.Name,
			&c.Year,
			&c.UniversityName,
		)
		if err != nil {
			http.Error(w, "Scan error", 500)
			return
		}

		// 🔥 Fetch departments for each course
		deptRows, err := database.DB.Query(`
			SELECT d.name
			FROM course_department cd
			JOIN department d ON cd.department_id = d.department_id
			WHERE cd.course_id = ?
		`, c.CourseID)

		if err != nil {
			http.Error(w, "Dept fetch error", 500)
			return
		}

		var departments []string
		for deptRows.Next() {
			var name string
			deptRows.Scan(&name)
			departments = append(departments, name)
		}
		deptRows.Close()

		c.Departments = departments

		courses = append(courses, c)
	}

	json.NewEncoder(w).Encode(courses)
}
