/*
connect.go file is responsible for connecting
to the database
*/
package models
import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"

)
func ConnectDatabase() {
	var db *gorm.DB
	db,err:=gorm.Open(sqlite.Open("bookstore.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
		/*The panic built-in function stops normal
		execution of the current goroutine. */
	}
	db.AutoMigrate(&Book{})
	
}
