package authentication

type UserHTTPReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDBData struct {
	Email    string
	Password string
	Role     string
	Name     string
}

type PasswordChangeReq struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}
