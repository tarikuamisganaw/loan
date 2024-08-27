package routers

import (
	"loan/controllers"
	"loan/delivery/controllers"
	"loan/infrastructure"
	"loan/usecases"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database, userUsecase usecases.UserUsecase, loanUsecase usecases.LoanUsecase, jwtService infrastructure.JWTService) *gin.Engine {
	router := gin.Default()

	// Initialize user controller
	userController := controllers.NewUserController(userUsecase)

	// Initialize loan controller
	loanController := controllers.NewLoanController(loanUsecase)

	// Admin Middleware
	adminMiddleware := infrastructure.AdminMiddleware(jwtService)

	// User Authentication Middleware
	authMiddleware := infrastructure.AuthMiddleware(jwtService)

	// User Routes
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", userController.RegisterUser)
		userRoutes.GET("/verify-email", userController.VerifyEmail)
		userRoutes.POST("/login", userController.Login)
	}

	// Loan Management Routes
	loanRoutes := router.Group("/loans")
	loanRoutes.Use(authMiddleware)
	{
		loanRoutes.POST("", loanController.ApplyForLoan)      // Apply for Loan
		loanRoutes.GET("/:id", loanController.ViewLoanStatus) // View Loan Status
	}

	// Admin Loan Management Routes
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(authMiddleware, adminMiddleware)
	{
		adminRoutes.GET("/loans", loanController.ViewAllLoans)                  // View All Loans
		adminRoutes.PATCH("/loans/:id/status", loanController.UpdateLoanStatus) // Approve/Reject Loan
		adminRoutes.DELETE("/loans/:id", loanController.DeleteLoan)             // Delete Loan
		adminRoutes.GET("/getallusers", userController.GetAllUsers)             // Get All Users (Admin)
		adminRoutes.DELETE("/deleteuser/:id", userController.DeleteUser)        // Delete User (Admin)
	}

	return router
}
