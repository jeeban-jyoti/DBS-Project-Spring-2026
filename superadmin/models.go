package superadmin

type AddEmployeeReq struct {
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         string  `json:"email"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone"`
	Gender        string  `json:"gender"`
	Salary        float64 `json:"salary"`
	AadhaarNumber string  `json:"aadhaar_number"`
}

type EmployeeDetail struct {
	EmployeeID    int     `json:"employee_id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         string  `json:"email"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone"`
	Gender        string  `json:"gender"`
	Salary        float64 `json:"salary"`
	AadhaarNumber string  `json:"aadhaar_number"`
}
