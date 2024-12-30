// user.route.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/controllers"
	"github.com/mmycin/ndc14/middlewares"
)

func SetupUserRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetUsers)

	group.POST("/signup", controllers.SignUp)
	group.POST("/login", controllers.Login)
	group.GET("/validate", middlewares.RequireAuth, controllers.Validate)
	group.GET("/logout", controllers.Logout)

	group.PUT("/update", middlewares.RequireAuth, controllers.UpdateUser)
	group.DELETE("/delete", middlewares.RequireAuth, controllers.DeleteUser)
}
