package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ConfigCors(router *gin.Engine) *gin.Engine {
	config := cors.DefaultConfig()

	// Allow all origins, or explicitly set the allowed origins.
	config.AllowAllOrigins = true // Allow any origin
	// config.AllowOrigins = []string{"http://localhost:4000"} // Allow only your frontend

	// Allow necessary methods and headers
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}

	// Apply CORS middleware
	router.Use(cors.New(config))

	// Return the router with CORS enabled
	return router
}
