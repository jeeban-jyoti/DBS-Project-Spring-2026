package admin

import (
	"encoding/json"
	"net/http"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func AddUniversity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddUniversityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec(`
		INSERT INTO university 
		(name, address, rep_first_name, rep_last_name, rep_email, rep_phone)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.Name, req.Address, req.RepFirstName, req.RepLastName, req.RepEmail, req.RepPhone)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "University added successfully",
		"university_id": id,
	})
}

func RemoveUniversity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "University ID required", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec(`
		DELETE FROM university WHERE university_id = ?
	`, id)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "University not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("University removed successfully"))
}
func UpdateUniversity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateUniversityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Build dynamic query
	query := "UPDATE university SET "
	args := []interface{}{}

	if req.Name != nil {
		query += "name = ?, "
		args = append(args, *req.Name)
	}
	if req.Address != nil {
		query += "address = ?, "
		args = append(args, *req.Address)
	}
	if req.RepFirstName != nil {
		query += "rep_first_name = ?, "
		args = append(args, *req.RepFirstName)
	}
	if req.RepLastName != nil {
		query += "rep_last_name = ?, "
		args = append(args, *req.RepLastName)
	}
	if req.RepEmail != nil {
		query += "rep_email = ?, "
		args = append(args, *req.RepEmail)
	}
	if req.RepPhone != nil {
		query += "rep_phone = ?, "
		args = append(args, *req.RepPhone)
	}

	// Remove trailing comma
	if len(args) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query = query[:len(query)-2] // remove ", "
	query += " WHERE university_id = ?"
	args = append(args, req.UniversityID)

	res, err := database.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "University not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "University updated successfully",
	})
}
