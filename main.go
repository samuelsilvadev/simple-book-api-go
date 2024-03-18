package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	db "simplebookapi/db"
	"simplebookapi/models"

	"github.com/gin-gonic/gin"
)

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, db.Books)
}

func createBook(context *gin.Context) {
	var newBook models.Book
	var err = context.BindJSON(&newBook)

	if err != nil {
		fmt.Println("Failed to parse request body.", err)
		return
	}

	if newBook.ID == 0 {
		lastBook := db.Books[len(db.Books)-1]
		newBook.ID = lastBook.ID + 1
	}

	hasBookByName := func() bool {
		for _, book := range db.Books {
			if book.Title == newBook.Title {
				return true
			}
		}

		return false
	}

	if hasBookByName() {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Book already exists"})
		return
	}

	db.Books = append(db.Books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

func checkOutBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}

	parsedId, _ := strconv.Atoi(id)

	findBookById := func() (models.Book, int, error) {
		for index, book := range db.Books {
			if book.ID == parsedId {
				return book, index, nil
			}
		}

		return models.Book{}, -1, errors.New("book not found")
	}

	book, bookIndex, err := findBookById()

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Book is not available"})
		return
	}

	book.Quantity = book.Quantity - 1
	db.Books[bookIndex] = book
	context.IndentedJSON(http.StatusOK, book)
}

func main() {
	fmt.Println("App is running...")

	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.PATCH("/books/checkout", checkOutBook)
	router.Run(":8888")
}
