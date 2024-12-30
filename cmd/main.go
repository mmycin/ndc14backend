package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/config"
	"github.com/mmycin/ndc14/libs"
	"github.com/mmycin/ndc14/middlewares"
	"github.com/mmycin/ndc14/routes"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
	config.SyncDatabase()
}

func main() {
	router := gin.Default()

	// Configure CORS
	router = middlewares.ConfigCors(router)

	// Create an API group
	api := router.Group("/api/v2")

	// Setup routes under the /api prefix
	routes.SetupUserRoutes(api.Group("/users"))
	routes.SetupNoticeRoutes(api.Group("/notices"))
	routes.SetupContactRoutes(api.Group("/contacts"))

	// Start server
	libs.Success("Starting server at " + os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	router.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
}
