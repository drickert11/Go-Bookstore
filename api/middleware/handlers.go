package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-bookstore/models"
	"log"
	"net/http"
	"os"      // used to read the environment variable
	"strconv" // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

//JSON structure
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	resultID := insertBook(book)

	response := response{
		ID:      resultID,
		Message: "Book was created seccessfully",
	}

	json.NewEncoder(w).Encode(response)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	book, err := getBook(int64(id))

	if err != nil {
		log.Fatalf("Unable to get book. %v", err)
	}

	json.NewEncoder(w).Encode(book)
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	books, err := getAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all books. %v", err)
	}

	json.NewEncoder(w).Encode(books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var book models.Book

	err = json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updateCount := updateBook(int64(id), book)

	msg := fmt.Sprintf("Book was updated successfully. Rowcount: %v", updateCount)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deleteCount := deleteBook(int64(id))

	msg := fmt.Sprintf("Book was deleted successfully. Rowcount: %v", deleteCount)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// ----------------- Data Access --------------------
func insertBook(book models.Book) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO books (title, author, publisher, publishDate, rating, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id` //stored proc

	var id int64

	err := db.QueryRow(sqlStatement, book.Title, book.Author, book.Publisher, book.PublishDate, book.Rating, book.Status).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

func getBook(id int64) (models.Book, error) {
	db := createConnection()

	defer db.Close()

	var book models.Book

	sqlStatement := `SELECT * FROM books WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher, &book.PublishDate, &book.Rating, &book.Status)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return book, nil
	case nil:
		return book, nil
	default:
		log.Fatalf("Issue with Row scan. %v", err)
	}

	return book, err
}

func getAllBooks() ([]models.Book, error) {
	db := createConnection()

	defer db.Close()

	var books []models.Book

	sqlStatement := `SELECT * FROM books`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var book models.Book

		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher, &book.PublishDate, &book.Rating, &book.Status)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		books = append(books, book)

	}

	return books, err
}

func updateBook(id int64, book models.Book) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `UPDATE books SET title=$2, author=$3, publisher=$4, publishdate=$5, rating=$6, status=$7 WHERE id=$1`

	res, err := db.Exec(sqlStatement, book.ID, book.Title, book.Author, book.Publisher, book.PublishDate, book.Rating, book.Status)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteBook(id int64) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM books WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
