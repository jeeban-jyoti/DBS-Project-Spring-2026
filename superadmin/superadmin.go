package superadmin

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

func AddAdmin(w http.ResponseWriter, r *http.Request) {
	var req AddEmployeeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	password := generatePassword(10)

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Insert into user
	res, err := tx.Exec(`
		INSERT INTO user (first_name, last_name, email, address, phone, password_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.FirstName, req.LastName, req.Email, req.Address, req.Phone, password)

	if err != nil {
		http.Error(w, "User insert failed", http.StatusInternalServerError)
		return
	}

	userID, _ := res.LastInsertId()

	// 2. Insert into employee
	_, err = tx.Exec(`
		INSERT INTO employee (employee_id, gender, salary, aadhaar_number)
		VALUES (?, ?, ?, ?)
	`, userID, req.Gender, req.Salary, req.AadhaarNumber)

	if err != nil {
		http.Error(w, "Employee insert failed", http.StatusInternalServerError)
		return
	}

	// 3. Insert into administrator
	_, err = tx.Exec(`
		INSERT INTO administrator (employee_id)
		VALUES (?)
	`, userID)

	if err != nil {
		http.Error(w, "Admin insert failed", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Admin added",
		"password": password,
	})
}

func RemoveAdmin(w http.ResponseWriter, r *http.Request) {
	empID := r.URL.Query().Get("empid")
	if empID == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	tx.Exec(`DELETE FROM administrator WHERE employee_id = ?`, empID)
	tx.Exec(`DELETE FROM employee WHERE employee_id = ?`, empID)
	tx.Exec(`DELETE FROM user WHERE user_id = ?`, empID)

	tx.Commit()

	w.Write([]byte("Admin removed successfully"))
}

func AddSupportStaff(w http.ResponseWriter, r *http.Request) {
	var req AddEmployeeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	password := generatePassword(10)

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	res, err := tx.Exec(`
		INSERT INTO user (first_name, last_name, email, address, phone, password_hash)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.FirstName, req.LastName, req.Email, req.Address, req.Phone, password)

	if err != nil {
		http.Error(w, "User insert failed", http.StatusInternalServerError)
		return
	}

	userID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		INSERT INTO employee (employee_id, gender, salary, aadhaar_number)
		VALUES (?, ?, ?, ?)
	`, userID, req.Gender, req.Salary, req.AadhaarNumber)

	if err != nil {
		http.Error(w, "Employee insert failed", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO customer_support (employee_id)
		VALUES (?)
	`, userID)

	if err != nil {
		http.Error(w, "Support insert failed", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Support staff added",
		"password": password,
	})
}

func RemoveSupportStaff(w http.ResponseWriter, r *http.Request) {
	empID := r.URL.Query().Get("empid")

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	tx.Exec(`DELETE FROM customer_support WHERE employee_id = ?`, empID)
	tx.Exec(`DELETE FROM employee WHERE employee_id = ?`, empID)
	tx.Exec(`DELETE FROM user WHERE user_id = ?`, empID)

	tx.Commit()

	w.Write([]byte("Support staff removed successfully"))
}

func FetchAdmins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT 
			u.user_id,
			u.first_name,
			u.last_name,
			u.email,
			u.address,
			u.phone,
			e.gender,
			e.salary,
			e.aadhaar_number
		FROM administrator a
		JOIN employee e ON a.employee_id = e.employee_id
		JOIN user u ON u.user_id = e.employee_id
	`)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var admins []EmployeeDetail

	for rows.Next() {
		var emp EmployeeDetail
		err := rows.Scan(
			&emp.EmployeeID,
			&emp.FirstName,
			&emp.LastName,
			&emp.Email,
			&emp.Address,
			&emp.Phone,
			&emp.Gender,
			&emp.Salary,
			&emp.AadhaarNumber,
		)
		if err != nil {
			http.Error(w, "Scan error", 500)
			return
		}
		admins = append(admins, emp)
	}

	json.NewEncoder(w).Encode(admins)
}

func FetchSupportStaff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT 
			u.user_id,
			u.first_name,
			u.last_name,
			u.email,
			u.address,
			u.phone,
			e.gender,
			e.salary,
			e.aadhaar_number
		FROM customer_support cs
		JOIN employee e ON cs.employee_id = e.employee_id
		JOIN user u ON u.user_id = e.employee_id
	`)
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	defer rows.Close()

	var staff []EmployeeDetail

	for rows.Next() {
		var emp EmployeeDetail
		err := rows.Scan(
			&emp.EmployeeID,
			&emp.FirstName,
			&emp.LastName,
			&emp.Email,
			&emp.Address,
			&emp.Phone,
			&emp.Gender,
			&emp.Salary,
			&emp.AadhaarNumber,
		)
		if err != nil {
			http.Error(w, "Scan error", 500)
			return
		}
		staff = append(staff, emp)
	}

	json.NewEncoder(w).Encode(staff)
}
