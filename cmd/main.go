package main

import (
	"log"

	"auth-barniee/internal/config"
	"auth-barniee/internal/database"
	"auth-barniee/internal/routes"

	"github.com/gin-gonic/gin"

	// Swagger imports
	_ "auth-barniee/cmd/docs" // Import generated docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Barniee Auth Service API
// @version 1.0
// @description This is the API documentation for Barniee Auth Service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@barniee.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.LoadConfig()

	db := database.InitDB(cfg)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupAuthRoutes(r, db, cfg)

	log.Printf("Auth service listening on port %s...", "8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Auth service failed to start: %v", err)
	}
}
