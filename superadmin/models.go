package superadmin

type AddUserReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
