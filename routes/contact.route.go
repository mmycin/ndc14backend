package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mmycin/ndc14/controllers"
)

func SetupContactRoutes(group *gin.RouterGroup) {
	group.GET("/contacts", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Contact route"})
	})
}
