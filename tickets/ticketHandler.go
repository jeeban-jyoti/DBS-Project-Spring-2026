package tickets

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jeeban-jyoti/DSB-Project-Spring-2026/database"
)

// ticket status (“new”, “assigned”, “in-process”, “completed”)

func GenerateTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTicketReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user := r.Header.Get("user_email")

	res, err := database.DB.Exec(
		`INSERT INTO tickets (category, title, problem_description, created_by, status)
		 VALUES (?, ?, ?, ?, 'new')`,
		req.Category, req.Title, req.ProblemDescription, user,
	)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "ticket created",
		"ticket_id": id,
	})
}

func ViewTickets(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, category, title, problem_description,
		solution_description, status, created_by, resolved_by, created_at, completed_at
		FROM tickets WHERE 1=1`

	var args []interface{}

	// ---- MULTIPLE STATUS SUPPORT ----
	statusParams := r.URL.Query()["status"] // handles ?status=a&status=b

	if len(statusParams) == 1 && strings.Contains(statusParams[0], ",") {
		// handle ?status=a,b,c
		statusParams = strings.Split(statusParams[0], ",")
	}

	if len(statusParams) > 0 {
		placeholders := make([]string, len(statusParams))
		for i, s := range statusParams {
			placeholders[i] = "?"
			args = append(args, s)
		}
		query += " AND status IN (" + strings.Join(placeholders, ",") + ")"
	}

	// ---- OTHER FILTERS ----
	createdBy := r.URL.Query().Get("created_by")
	resolvedBy := r.URL.Query().Get("resolved_by")

	if createdBy != "" {
		query += " AND created_by = ?"
		args = append(args, createdBy)
	}

	if resolvedBy != "" {
		query += " AND resolved_by = ?"
		args = append(args, resolvedBy)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tickets []Ticket

	for rows.Next() {
		var t Ticket
		err := rows.Scan(
			&t.ID, &t.Category, &t.Title, &t.ProblemDescription,
			&t.SolutionDescription, &t.Status, &t.CreatedBy,
			&t.ResolvedBy, &t.CreatedAt, &t.CompletedAt,
		)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		tickets = append(tickets, t)
	}

	json.NewEncoder(w).Encode(tickets)
}

func HandleNewTickets(w http.ResponseWriter, r *http.Request) {
	handleStatusUpdate(w, r, "new", "assigned")
}

func HandleAssignedTickets(w http.ResponseWriter, r *http.Request) {
	handleStatusUpdate(w, r, "assigned", "in-process")
}

func HandleInProcessTickets(w http.ResponseWriter, r *http.Request) {
	handleStatusUpdate(w, r, "in-process", "completed")
}

func handleStatusUpdate(w http.ResponseWriter, r *http.Request, from, to string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user := r.Header.Get("user_email")

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}

	var current string
	err = tx.QueryRow("SELECT status FROM tickets WHERE id = ?", id).Scan(&current)
	if err == sql.ErrNoRows {
		tx.Rollback()
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}
	if err != nil {
		tx.Rollback()
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if current != from {
		tx.Rollback()
		http.Error(w, "Invalid state transition", http.StatusBadRequest)
		return
	}

	_, err = tx.Exec("UPDATE tickets SET status = ? WHERE id = ?", to, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`INSERT INTO ticket_status_history 
		(ticket_id, old_status, new_status, changed_by)
		VALUES (?, ?, ?, ?)`,
		id, from, to, user,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "History insert failed", http.StatusInternalServerError)
		return
	}

	if from == "new" && to == "assigned" {
		_, err = tx.Exec(
			"UPDATE tickets SET status = ?, assigned_by = ? WHERE id = ?",
			to, user, id,
		)
	} else if from == "assigned" && to == "in-process" {
		_, err = tx.Exec(
			"UPDATE tickets SET status = ?, in_process_by = ? WHERE id = ?",
			to, user, id,
		)
	} else if to == "completed" {
		_, err = tx.Exec(
			`UPDATE tickets 
			 SET status = ?, resolved_by = ?, completed_at = NOW() 
			 WHERE id = ?`,
			to, user, id,
		)
	} else {
		_, err = tx.Exec(
			"UPDATE tickets SET status = ? WHERE id = ?",
			to, id,
		)
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	if err = tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "status updated to " + to,
	})
}
