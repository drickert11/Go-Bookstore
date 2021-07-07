package main

import (
	"fmt"
	"go-bookstore/middleware"
	"go-bookstore/router"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	r := router.Router()
	middleware.ResetDB()
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))

}
