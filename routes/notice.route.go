package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/controllers"
	"github.com/mmycin/ndc14/middlewares"
)

func SetupNoticeRoutes(group *gin.RouterGroup) {
	group.GET("/", controllers.GetNotices)
	group.GET("/:id", middlewares.RequireAuth, controllers.GetNotice)
	group.POST("/", middlewares.RequireAuth, controllers.CreateNotice)
	group.PUT("/:id", middlewares.RequireAuth, controllers.UpdateNotice)
	group.DELETE("/:id", middlewares.RequireAuth, controllers.DeleteNotice)
}
