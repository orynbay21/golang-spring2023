// menu file
package main

import (
	"github.com/gin-gonic/gin"
	//"github.com/pp221b030707/golang-project/controller"
	"github.com/pp221b030707/golang-project/models"

)

func main() {
	r := gin.Default()
	models.ConnectDatabase()
	//r.GET("/books", controller.FindBooks)
	//r.GET("/books/:id", controller.FindBook)
	//r.POST("/books", controller.CreateBook)
	//r.PATCH("/books/:id", controller.UpdateBook)
	//r.DELETE("/books/:id", controller.DeleteBook)

	r.Run()
}
