package main

import (
	"errors"
	"net/http"

	//"errors"
	"github.com/gin-gonic/gin"
)

// Table for the books
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// Collection of Library books (example)
var books = []book{
	{ID: "1", Title: "Little Women", Author: "Louisa May Alcott", Quantity: 3},
	{ID: "2", Title: "Pride, Price and Prejudice", Author: "Jane Austen", Quantity: 1},
	{ID: "3", Title: "Red Rising", Author: "Pierce Brown", Quantity: 4},
	{ID: "4", Title: "The Master and Margarita", Author: "Mikhail Bulgakov", Quantity: 2},
}

// Creating Get Request (Implementation of the API)
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

// Checking in and out of Library with books. Quantity changes with vice versa actions
// Checking out book by the ID as a query param.
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book unavailable."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

// Return book feature
// Same query parameter, so one can return a book via the ID No.
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing ID query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// New feature: Finding a book by ID (Checking in/out, getting a book by ID)
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

// Creating Post Request
// Created variable (newBook which is under type book)
// Attempted to binds JSON of request data to the new book by passing it's pointer.
// This method allows the error to be handled and get a response (.BindJSON)
// The return command sends back an error response message back to the sender
//Below, if there was no error we were successful with binding JSON to the struct then we have established a new book
// The book was created by simply appending it to the books slice.
// Then where is noted indented: we return the book create with the stat. code

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
