package routers

import (
	"loan/delivery/controllers"
	"loan/infrastructure"
	"loan/usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database, userUsecase usecases.UserUsecase, jwtService infrastructure.JWTService) *gin.Engine {
	router := gin.Default()

	// Initialize repository

	// Initialize controller
	userController := controllers.NewUserController(userUsecase)

	adminMiddleware := infrastructure.AdminMiddleware(jwtService)

	// User Routes
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", userController.RegisterUser)
		userRoutes.GET("/verify-email", userController.VerifyEmail)
		userRoutes.POST("/login", userController.Login)
	}
	auth := router.Group("/api")
	auth.Use(infrastructure.AuthMiddleware(jwtService))
	{
		auth.GET("/getallusers", adminMiddleware, userController.GetAllUsers)
		auth.DELETE("/deleteuser/:id", adminMiddleware, userController.DeleteUser)
	}

	return router
}
