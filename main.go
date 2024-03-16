package main

import (
	"fmt"
	"net/http"

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

func main() {
	fmt.Println("App is running...")

	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.Run(":8888")
}
