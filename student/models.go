package student

type Book struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Authors         string  `json:"authors"`
	PublicationDate string  `json:"publication_date"`
	Price           float64 `json:"price"`
	Quantity        int     `json:"quantity"`
}

type Review struct {
	StudentID int    `json:"student_id"`
	Rating    int    `json:"rating"`
	Text      string `json:"review_text"`
	Date      string `json:"review_date"`
}

type BookDetail struct {
	Title           string  `json:"title"`
	ISBN            string  `json:"isbn"`
	Publisher       string  `json:"publisher"`
	PublicationDate string  `json:"publication_date"`
	Edition         string  `json:"edition"`
	Language        string  `json:"language"`
	Format          string  `json:"format"`
	Type            string  `json:"type"`
	PurchaseOption  string  `json:"purchase_option"`
	Price           float64 `json:"price"`
	Quantity        int     `json:"quantity"`

	Authors       string  `json:"authors"`
	SubCategories string  `json:"subcategories"`
	AvgRating     float64 `json:"avg_rating"`

	Reviews []Review `json:"reviews"`
}

type BookFullResponse struct {
	BookID          int      `json:"id"`
	Title           string   `json:"title"`
	ISBN            string   `json:"isbn"`
	Publisher       string   `json:"publisher"`
	PublicationDate string   `json:"publication_date"`
	Edition         string   `json:"edition"`
	Language        string   `json:"language"`
	Format          string   `json:"format"`
	Type            string   `json:"type"`
	PurchaseOption  string   `json:"purchase_option"`
	Price           float64  `json:"price"`
	Quantity        int      `json:"quantity"`
	Category        string   `json:"category"`
	Subcategories   []string `json:"subcategories"`
	Authors         []string `json:"authors"`
	Keywords        []string `json:"keywords"`
}

type AddToCartReq struct {
	BookID   int `json:"bookid"`
	Quantity int `json:"quantity"`
}

type RemoveFromCartReq struct {
	BookID int `json:"bookid"`
}
