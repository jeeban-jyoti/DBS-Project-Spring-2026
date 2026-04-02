package tickets

type CreateTicketReq struct {
	Category           string `json:"category"`
	Title              string `json:"title"`
	ProblemDescription string `json:"problem_description"`
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
