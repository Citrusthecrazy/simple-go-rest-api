package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	//"errors"
)

type Book struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func getBooks(ctx *gin.Context) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database!")
	}

	var results []Book
	db.Table("books").Find(&results)
	ctx.IndentedJSON(http.StatusOK, results)
}

func getBookById(id int) (*Book, error) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")

	}
	var book Book
	db.Table("books").First(&book, id)
	if err := db.Where("ID = ?", id).First(&book).Error; err != nil {
		return nil, errors.New("Book not found")
	}

	return &book, nil

}

func createBook(ctx *gin.Context) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database!")
	}
	var newBook Book

	if err := ctx.BindJSON(&newBook); err != nil {
		return
	}

	db.Create(&Book{ID: newBook.ID, Title: newBook.Title, Author: newBook.Author, Quantity: newBook.Quantity})
	ctx.IndentedJSON(http.StatusCreated, newBook)
}

func bookById(ctx *gin.Context) {

	id := ctx.Param("id")
	i, parseErr := strconv.Atoi(id)

	if parseErr != nil {
		panic("There was a parsing error")
	}

	book, error := getBookById(i)

	if error != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(ctx *gin.Context) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database!")
	}

	id, ok := ctx.GetQuery("id")
	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query params"})
		return
	}

	i, parseErr := strconv.Atoi(id)
	if parseErr != nil {
		panic("There was a parsing error")
	}
	book, err := getBookById(i)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity -= 1
	db.Save(&book)
	ctx.IndentedJSON(http.StatusOK, book)
}

func returnBook(ctx *gin.Context) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database!")
	}
	id, ok := ctx.GetQuery("id")
	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query params"})
		return
	}
	i, parseErr := strconv.Atoi(id)
	if parseErr != nil {
		panic("There was a parsing error")
	}
	book, err := getBookById(i)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	book.Quantity += 1
	db.Save(&book)
	ctx.IndentedJSON(http.StatusOK, book)

}

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database!")
	}

	db.AutoMigrate(&Book{})
	router := gin.Default()

	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)

	router.Run("localhost:8080")
}
