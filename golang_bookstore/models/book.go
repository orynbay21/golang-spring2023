package models
import(
	"gorm.io/gorm"
	//"gorm.io/driver/sqlite"	
)
type Book struct {
	gorm.Model
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Description string `json:"description"`
	Cost uint `json:"cost"`
}
/*
model- obj that represents
 the data in ur application
In our case its a -book


*/