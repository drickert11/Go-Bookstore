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
