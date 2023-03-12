package main

import (
	"github.com/gin-gonic/gin"
	
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	/* routes */
	r.GET("/books", controller.FindBooks)
	r.GET("/books/:id", controller.FindBook)
	r.POST("/books", controller.CreateBook)
	r.PATCH("/books/:id", controller.UpdateBook)
	r.DELETE("/books/:id", controller.DeleteBook)

	r.Run()
	/*
	
	get book information (title, description, cost) by ID
	Get list of all books
	Update title, description of book by ID
	Delete a book by ID
	Search a book by title
	Add a book
	Get a sorted list of books ordered by cost in descending, ascending orders
	For database you have to use Gorm. Tables should be normalized.
	Write dockerfile which runs all the above-mentioned endpoints.
	Build an image based on this dockerfile. Launch container based on the created image.

	*/
}
