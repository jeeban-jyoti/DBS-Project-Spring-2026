package admin

// University
type AddUniversityReq struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	RepFirstName string `json:"rep_first_name"`
	RepLastName  string `json:"rep_last_name"`
	RepEmail     string `json:"rep_email"`
	RepPhone     string `json:"rep_phone"`
}

type UpdateUniversityReq struct {
	UniversityID int     `json:"university_id"`
	Name         *string `json:"name"`
	Address      *string `json:"address"`
	RepFirstName *string `json:"rep_first_name"`
	RepLastName  *string `json:"rep_last_name"`
	RepEmail     *string `json:"rep_email"`
	RepPhone     *string `json:"rep_phone"`
}

// Book
type AddBookReq struct {
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

	Category      string   `json:"category"`
	Subcategories []string `json:"subcategories"`
	Authors       []string `json:"authors"`
	Keywords      []string `json:"keywords"`
}

//Department

type AddDepartmentReq struct {
	Name         string `json:"name"`
	UniversityID int    `json:"university_id"`
}

type CourseReq struct {
	Name         string `json:"name"`
	UniversityID int    `json:"university_id"`
	Year         int    `json:"year"`
	Departments  []int  `json:"departments"`
}

type InstructorReq struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	UniversityID int    `json:"university_id"`
	DepartmentID int    `json:"department_id"`
}

type StudentReq struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	DOB          string `json:"date_of_birth"`
	UniversityID int    `json:"university_id"`
	Major        string `json:"major"`
	Status       string `json:"status"`
	YearOfStudy  int    `json:"year_of_study"`
}

// Semester
type AddSemesterReq struct {
	Year         int    `json:"year"`
	Season       string `json:"season"`
	CourseID     int    `json:"course_id"`
	InstructorID int    `json:"instructor_id"`
	UniversityID int    `json:"university_id"`
	BookIDs      []int  `json:"book_ids"` // optional
}

type SemesterDetail struct {
	SemID          int      `json:"sem_id"`
	Year           int      `json:"year"`
	Season         string   `json:"season"`
	CourseName     string   `json:"course_name"`
	InstructorName string   `json:"instructor_name"`
	UniversityName string   `json:"university_name"`
	Books          []string `json:"books"`
}

type CourseDetail struct {
	CourseID       int      `json:"course_id"`
	Name           string   `json:"name"`
	Year           int      `json:"year"`
	UniversityName string   `json:"university_name"`
	Departments    []string `json:"departments"`
}
