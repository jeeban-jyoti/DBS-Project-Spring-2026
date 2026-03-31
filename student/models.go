package student

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	PubYear  int    `json:"pub_year"`
	Status   string `json:"status"`
	Quantity int    `json:"quantity"`
}
