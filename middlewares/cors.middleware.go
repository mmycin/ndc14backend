package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ConfigCors(router *gin.Engine) *gin.Engine {
	config := cors.DefaultConfig()

	// Allow specific origins for security, or allow all origins for development
	config.AllowOrigins = []string{"http://localhost:4000", "http://your-production-frontend.com"} // Change to match your frontend URLs
	config.AllowCredentials = true // Required for cookies and authentication headers

	// Allow necessary methods
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE", "PATCH"}

	// Allow required headers
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length", "Authorization"} // Expose headers to the client

	// Apply CORS middleware
	router.Use(cors.New(config))

	return router
}
