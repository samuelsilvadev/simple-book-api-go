package main

import (
	"fmt"
	"net/http"

	db "simplebookapi/db"

	"github.com/gin-gonic/gin"
)

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, db.Books)
}

func main() {
	fmt.Println("App is running...")

	router := gin.Default()
	router.GET("/books", getBooks)
	router.Run(":8888")
}
