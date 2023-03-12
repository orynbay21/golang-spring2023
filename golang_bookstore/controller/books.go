package controller

import (
	//"net/http"

	//"github.com/gin-gonic/gin"
	//"github.com/pp221b030707/golang-project/models"
	"gorm.io/gorm"
)

type CreateBookInput struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Cost        uint   `json:"cost"`
}
type UpdateBookInput struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Cost        uint   `json:"cost"`
}
