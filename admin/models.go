package admin

type Student struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AddStudentsReq struct {
	Students []Student `json:"students"`
}

type RemoveStudentsReq struct {
	Emails []string `json:"emails"`
}

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	PubYear  int    `json:"pub_year"`
	Status   string `json:"status"`
	Quantity int    `json:"quantity"`
}

type AddBooksReq struct {
	Books []Book `json:"books"`
}

type RemoveBooksReq struct {
	IDs []int `json:"ids"`
}
