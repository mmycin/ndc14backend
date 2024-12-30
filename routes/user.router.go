package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/controllers"
	"github.com/mmycin/ndc14/middlewares"
)

func SetupUserRoutes(group *gin.RouterGroup) {
	group.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "User route"})
	})

	group.POST("/signup", controllers.SignUp)
	group.POST("/login", controllers.Login)
	group.GET("/validate", middlewares.RequireAuth, controllers.Validate)
	group.GET("/logout", controllers.Logout)
}
