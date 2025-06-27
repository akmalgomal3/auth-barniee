package routes

import (
	"auth-barniee/internal/config"
	"auth-barniee/internal/handlers"
	"auth-barniee/internal/middlewares"
	"auth-barniee/internal/repositories"
	"auth-barniee/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	configCors.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	configCors.ExposeHeaders = []string{"Content-Length"}
	configCors.AllowCredentials = true
	configCors.MaxAge = 30

	r.Use(cors.New(configCors))

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	schoolRepo := repositories.NewSchoolRepository(db)
	packageRepo := repositories.NewPackageRepository(db)
	emailVerifyRepo := repositories.NewEmailVerificationRepository(db)

	authService := services.NewAuthService(userRepo, roleRepo, schoolRepo, cfg) // Updated
	userService := services.NewUserService(userRepo, roleRepo)
	registrationService := services.NewRegistrationService(schoolRepo, userRepo, roleRepo, packageRepo, emailVerifyRepo, cfg)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	registrationHandler := handlers.NewRegistrationHandler(registrationService)

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)

		registration := public.Group("/register")
		{
			registration.POST("/school-info", registrationHandler.RegisterSchoolInfo)
			registration.POST("/admin-info", registrationHandler.RegisterAdminInfo)
			registration.GET("/packages", registrationHandler.GetAllPackages)
			registration.POST("/select-package", registrationHandler.SelectPackage)
			registration.POST("/email-verification/request-otp", registrationHandler.RequestEmailVerificationOTP)
			registration.POST("/email-verification/verify-otp", registrationHandler.VerifyEmailOTP)
			registration.POST("/complete", registrationHandler.CompleteRegistration)
		}
	}

	authenticated := r.Group("/api/v1")
	authenticated.Use(middlewares.AuthMiddleware(cfg))
	{
		authenticated.POST("/auth/logout", authHandler.Logout)

		authenticated.GET("/profile", authHandler.GetUserProfile)

		admin := authenticated.Group("/admin")
		admin.Use(middlewares.AuthorizeRoles("admin"))
		{
			admin.POST("/users", userHandler.CreateTeacherOrStudent)
			admin.GET("/users", userHandler.GetAllUsers)
			admin.GET("/users/:id", userHandler.GetUserByID)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
		}
	}
}
