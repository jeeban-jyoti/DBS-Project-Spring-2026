package authentication

type UserHTTPReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDBData struct {
	email    string
	password string
	role     string
	name     string
}
