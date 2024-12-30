package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/controllers"
	_ "github.com/mmycin/ndc14/controllers"
)

func SetupNoticeRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetNotices)
}
