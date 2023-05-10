package routes

import (
	"golangfinal/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRoutes.POST("/admin/addrating", controllers.AddRating())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
	incomingRoutes.GET("/users/filterequalprice", controllers.FilterEqualPrice())
	incomingRoutes.GET("/users/filterlowerprice", controllers.FilterLowerPrice())

}
