package student

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	PubYear  int    `json:"pub_year"`
	Status   string `json:"status"`
	Quantity int    `json:"quantity"`
}

type AddToCartReq struct {
	UserID   string `json:"userid"`
	BookID   int    `json:"bookid"`
	Quantity int    `json:"quantity"`
}

type RemoveFromCartReq struct {
	UserID string `json:"userid"`
	BookID int    `json:"bookid"`
}

type OrderPlaceReq struct {
	UserID string `json:"userid"`
}
