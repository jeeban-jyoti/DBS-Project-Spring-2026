package tickets

import "database/sql"

type GenerateNewTicketReq struct {
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AssignTicketReq struct {
	TicketID int `json:"ticket_id"`
	AdminID  int `json:"admin_id"`
}

type ChangeTicketSStatusReq struct {
	TicketID            int    `json:"ticket_id"`
	NewStatus           string `json:"new_status"`
	SolutionDescription string `json:"solution_description"`
}

type Ticket struct {
	ID                  int     `json:"id"`
	Category            string  `json:"category"`
	Title               string  `json:"title"`
	ProblemDescription  string  `json:"problem_description"`
	SolutionDescription *string `json:"solution_description"`
	Status              string  `json:"status"`
	CreatedBy           string  `json:"created_by"`
	ResolvedBy          *string `json:"resolved_by"`
	CreatedAt           string  `json:"created_at"`
	CompletedAt         *string `json:"completed_at"`
}

type TicketResponse struct {
	TicketID            int            `json:"ticket_id"`
	GeneratedBy         string         `json:"generated_by"` // User's full name or email
	Category            string         `json:"category"`
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	SolutionDescription sql.NullString `json:"solution_description"`
	CreatedDate         string         `json:"created_date"`
	CompletionDate      sql.NullString `json:"completion_date"`
	Status              string         `json:"status"`
	AssignedAdminID     sql.NullInt64  `json:"assigned_admin_id"`
}
