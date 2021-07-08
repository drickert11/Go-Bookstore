package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-bookstore/models"
	"log"
	"net/http" // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	// package used to read the .env file
	_ "github.com/lib/pq" // postgres golang driver
)

//JSON structure
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

var connectionString = "postgres://aprnxuwd:sqWwoqHvVOvEgnPBOLJHb8jqiBzuL6xu@batyr.db.elephantsql.com/aprnxuwd"

func createConnection() *sql.DB {

	db, err := sql.Open("postgres", connectionString)

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
	} else if check := models.IsValidBook(book); check != "none" {
		respondWithError(w, http.StatusNotAcceptable, fmt.Sprintf("Book was invalid because of %v", check))
		return
	}

	resultID := insertBook(book)

	response := response{
		ID:      resultID,
		Message: "Book was created seccessfully",
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	log.Printf("Params are %v", params)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	book, err := getBook(int64(id))
	fmt.Printf("error says %v", err)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("No entry found with the id= %v", id))
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Some problem occurred.  %v", err))
		return
	default:
		respondWithJSON(w, http.StatusOK, book)
	}
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	books, err := getAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all books. %v", err)
	}

	respondWithJSON(w, http.StatusOK, books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	} else if check := models.IsValidBook(book); check != "none" {
		respondWithError(w, http.StatusNotAcceptable, fmt.Sprintf("Book was invalid because of %v", check))
		return
	}

	updateCount := updateBook(book)

	msg := fmt.Sprintf("Book was updated successfully. Rowcount: %v", updateCount)

	res := response{
		ID:      int64(book.ID),
		Message: msg,
	}

	respondWithJSON(w, http.StatusOK, res)
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
	respondWithJSON(w, http.StatusOK, res)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ----------------- Data Access --------------------
const tableCreationQuery = `CREATE TABLE IF NOT EXISTS books
	(
		id INT GENERATED ALWAYS AS IDENTITY,
		title varchar(100) NOT NULL,
		author varchar(100) NOT NULL,
		publisher varchar(100) NOT NULL,
		publishDate date NOT NULL,
		rating int NOT NULL,
		status varchar(10) NOT NULL
	)`

func ResetDB() {
	db := createConnection()

	if _, err1 := db.Exec(tableCreationQuery); err1 != nil {
		log.Fatalf("Unable to execute table check. %v", err1)
	}

	defer db.Close()
	sqlstatement := `Truncate TABLE books RESTART IDENTITY`

	if _, err := db.Exec(sqlstatement); err != nil {
		log.Fatalf("Unable to execute the TRUNCATE command. %v", err)
	}

	var id int64

	sqlstatement2 := `INSERT INTO books(title, author, publisher, publishDate, rating, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID`
	err2 := db.QueryRow(sqlstatement2, "test", "test", "test", "1/1/2020", 1, "CheckedIn").Scan(&id)

	if err2 != nil {
		log.Fatalf("Unable to execute the INSERT command. %v", err2)
	}

	fmt.Printf("DB was reset and row with ID: %v was created", id)

	err3 := db.QueryRow(sqlstatement2, "test2", "test2", "test2", "1/2/2021", 3, "CheckedOut").Scan(&id)

	if err3 != nil {
		log.Fatalf("Unable to execute the INSERT command. %v", err3)
	}

	fmt.Printf("DB was reset and row with ID: %v was created", id)
}

func insertBook(book models.Book) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO books (title, author, publisher, publishDate, rating, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

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
		return book, sql.ErrNoRows
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

func updateBook(book models.Book) int64 {

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
