package router

import (
	"go-bookstore/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/book/{id}", middleware.GetBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/book", middleware.GetAllBooks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newbook", middleware.AddBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/book", middleware.UpdateBook).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletebook/{id}", middleware.DeleteBook).Methods("DELETE", "OPTIONS")

	return router
}
