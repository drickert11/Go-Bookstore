Change Directory to "Go-Bookstore"

To build, run this command in terminal:
docker build -t my-bookstore .

To run, run this command in terminal:
docker run -p 8080:8080 -tid my-bookstore

You can grab the Container ID with this command, the ID should pop up to the left of the name my-bookstore that was set when image was built:
docker ps

With the id you can stop the container with this command:
docker stop [container_id]

You can also remove a stopped container altogether with this command:
docker rm [container_id]

Request and Responses are handled in JSON format
The server runs via localhost on port 8080, and we have a few routes to choose:

http://localhost:8080/api/book will return all of the books in the db

http://localhost:8080/api/book/[id] (example: http://localhost:8080/api/book/1) will return a book through it's id, or an error if id is not valid

http://localhost:8080/api/newbook is a post method to create books with the
shape
{
    "id":0,
    "title":"Book3",
    "author":"John Smith",
    "publisher":"test publisher",
    "publishDate":"2003-01-23",
    "rating":3,
    "status":"CheckedIn"
} 
id can be left as 0 or anything when adding a book.

http://localhost:8080/api/book as a PUT request will accept the same shape as a new book, however the id does matter this time and will update the book at that id or throw an error.

http://localhost:8080/api/deletebook/[id] (example: http://localhost:8080/api/deletebook/3) is a DELETE request that will purge the book of that id from the database or throw an error if it can't find it.

For Testing:
Within the Go_Boostore directory you can run:
go test ./... 
To note: all of the tests were stored in Middleware.



