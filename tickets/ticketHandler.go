package tickets

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

func getUserID(email string) (int, error) {
	var id int
	err := database.DB.QueryRow(`
		SELECT user_id FROM user where email = ?
	`, email).Scan(&id)
	return id, err
}

// 1. GenerateNewTicket: POST /api/v1/generateTicket
func GenerateNewTicket(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "user_email header is required", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req GenerateNewTicketReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO ticket (created_by_user_id, category, title, description, status) 
		VALUES (?, ?, ?, ?, 'new')
	`, userID, req.Category, req.Title, req.Description)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket generated successfully"})
}

// 2. AssignTicket: POST /api/v1/assignTicket
func AssignTicket(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "user_email header is required", http.StatusBadRequest)
		return
	}

	supportStaffID, err := getUserID(email)
	if err != nil {
		http.Error(w, "Support staff user not found", http.StatusNotFound)
		return
	}

	var req AssignTicketReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Update Ticket
	_, err = tx.Exec(`
		UPDATE ticket 
		SET assigned_admin_id = ?, status = 'assigned' 
		WHERE ticket_id = ? AND status = 'new'
	`, req.AdminID, req.TicketID)

	if err != nil {
		http.Error(w, "Failed to assign ticket", http.StatusInternalServerError)
		return
	}

	// Log to History using the ID retrieved from email
	_, err = tx.Exec(`
		INSERT INTO ticket_status_history (ticket_id, changed_by, old_status, new_status) 
		VALUES (?, ?, 'new', 'assigned')
	`, req.TicketID, supportStaffID)

	if err != nil {
		http.Error(w, "Failed to log history", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket assigned successfully"})
}

// 3. ChangeTicketStatus: PUT /api/v1/changeTicketStatus
func ChangeTicketStatus(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "user_email header is required", http.StatusBadRequest)
		return
	}

	adminID, err := getUserID(email)
	if err != nil {
		http.Error(w, "Admin user not found", http.StatusNotFound)
		return
	}

	var req ChangeTicketSStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status enum
	validStatuses := map[string]bool{"assigned": true, "in-process": true, "completed": true}
	if !validStatuses[req.NewStatus] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Get current status and verify assignment using the ID retrieved from email
	var oldStatus string
	err = tx.QueryRow(`SELECT status FROM ticket WHERE ticket_id = ? AND assigned_admin_id = ?`, req.TicketID, adminID).Scan(&oldStatus)
	if err == sql.ErrNoRows {
		http.Error(w, "Ticket not found or not assigned to you", http.StatusForbidden)
		return
	}

	// 2. Update Ticket
	completionDate := sql.NullTime{}
	if req.NewStatus == "completed" {
		completionDate = sql.NullTime{Time: time.Now(), Valid: true}
	}

	_, err = tx.Exec(`
		UPDATE ticket 
		SET status = ?, solution_description = ?, completion_date = ? 
		WHERE ticket_id = ?
	`, req.NewStatus, req.SolutionDescription, completionDate, req.TicketID)

	if err != nil {
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		return
	}

	// 3. Log to History
	_, err = tx.Exec(`
		INSERT INTO ticket_status_history (ticket_id, changed_by, old_status, new_status) 
		VALUES (?, ?, ?, ?)
	`, req.TicketID, adminID, oldStatus, req.NewStatus)

	if err != nil {
		http.Error(w, "Failed to log history", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket status updated successfully"})
}

// Helper to map database rows to TicketResponse
func scanTicketRow(rows *sql.Rows) (TicketResponse, error) {
    var tr TicketResponse
    var createdDate []byte

    err := rows.Scan(
        &tr.TicketID,
        &tr.GeneratedBy,
		&tr.GeneratedByEmail,
        &tr.Category,
        &tr.Title,
        &tr.Description,
        &tr.SolutionDescription,
        &createdDate,
        &tr.CompletionDate,
        &tr.Status,
        &tr.AssignedAdminID,
    )
    if err != nil {
        return tr, err
    }
    createdDateParsed, err := time.Parse("2006-01-02 15:04:05", string(createdDate))
    if err != nil {
        return tr, err
    }
    tr.CreatedDate = createdDateParsed.Format("2006-01-02 15:04:05")
    return tr, nil
}

// 4. ShowGeneratedTickets: GET /api/v1/myTickets
func ShowGeneratedTickets(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "user_email header is required", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Println(userID)

	rows, err := database.DB.Query(`
		SELECT t.ticket_id, CONCAT(u.first_name, ' ', u.last_name), u.email, t.category, t.title, 
		       t.description, t.solution_description, t.created_date, 
		       t.completion_date, t.status, t.assigned_admin_id
		FROM ticket t
		JOIN user u ON t.created_by_user_id = u.user_id
		WHERE t.created_by_user_id = ?
	`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tickets []TicketResponse
	for rows.Next() {
		ticket, err := scanTicketRow(rows)
		if err != nil {
			fmt.Println("Scan error:", err)
			http.Error(w, "Error reading tickets", http.StatusInternalServerError)
			return
		}
		tickets = append(tickets, ticket)
	}

	json.NewEncoder(w).Encode(tickets)
}

// 5a. ShowNewTickets: GET /api/v1/viewNewTickets (Support Staff only — new tickets for assignment)
func ShowNewTickets(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "user_email header is required", http.StatusBadRequest)
		return
	}

	_, err := getUserID(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query(`
		SELECT t.ticket_id, CONCAT(u.first_name, ' ', u.last_name), u.email, t.category, t.title, 
		       t.description, t.solution_description, t.created_date, 
		       t.completion_date, t.status, t.assigned_admin_id
		FROM ticket t
		JOIN user u ON t.created_by_user_id = u.user_id
		WHERE t.status = 'new'
		ORDER BY t.created_date ASC
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tickets []TicketResponse
	for rows.Next() {
		ticket, err := scanTicketRow(rows)
		if err != nil {
			continue
		}
		tickets = append(tickets, ticket)
	}

	json.NewEncoder(w).Encode(tickets)
}

// 5b. ShowALLTickets: GET /api/v1/viewALLTickets (Admin only)
func ShowALLTickets(w http.ResponseWriter, r *http.Request) {
	// Even for admin, check if the user exists in the system
	email := r.Header.Get("user_email")
	if email == "" {
		http.Error(w, "user_email header is required", http.StatusBadRequest)
		return
	}

	_, err := getUserID(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query(`
		SELECT t.ticket_id, CONCAT(u.first_name, ' ', u.last_name), u.email, t.category, t.title, 
		       t.description, t.solution_description, t.created_date, 
		       t.completion_date, t.status, t.assigned_admin_id
		FROM ticket t
		JOIN user u ON t.created_by_user_id = u.user_id
		WHERE t.status != 'new'
		ORDER BY t.created_date DESC
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tickets []TicketResponse
	for rows.Next() {
		ticket, err := scanTicketRow(rows)
		if err != nil {
			continue
		}
		tickets = append(tickets, ticket)
	}

	json.NewEncoder(w).Encode(tickets)
}
