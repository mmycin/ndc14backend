package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mmycin/ndc14/controllers"
)

func SetupNoticeRoutes(group *gin.RouterGroup) {
	group.GET("/notices", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Notice route"})
	})
}
