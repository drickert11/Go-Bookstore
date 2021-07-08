package models

type Book struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publishDate"`
	Rating      int64  `json:"rating"`
	Status      string `json:"status"`
}

func IsValidBook(book Book) string {
	//TODO: replace with reflection
	switch {
	case book.Title == "":
		return "Title"
	case book.Author == "":
		return "Author"
	case book.Publisher == "":
		return "Publisher"
	case book.Rating < 1 || book.Rating > 3:
		return "Rating"
	case !(book.Status == "CheckedIn" || book.Status == "CheckedOut"):
		return "Status"
	default:
		return "none"
	}
}
