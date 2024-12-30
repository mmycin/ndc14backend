package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/controllers"
	"github.com/mmycin/ndc14/middlewares"
)

func SetupContactRoutes(group *gin.RouterGroup) {
	// Public route - anyone can create a contact
	group.POST("/", controllers.CreateContact)

	// Protected routes - require authentication
	group.GET("/", middlewares.RequireAuth, controllers.GetContacts)
	group.GET("/:id", middlewares.RequireAuth, controllers.GetContact)
	group.DELETE("/:id", middlewares.RequireAuth, controllers.DeleteContact)
}
