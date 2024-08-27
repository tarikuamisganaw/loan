package routers

import (
	"loan/delivery/controllers"
	"loan/domain"
	"loan/infrastructure"

	"github.com/gin-gonic/gin"
)

func NewRefreshTokenRouter(r *gin.Engine, userUsecase domain.UserUsecaseInterface, refreshTokenUsecase domain.RefreshTokenUsecaseInterface, jwtService infrastructure.JWTService) {
	refreshTokenController := controllers.NewRefreshTokenController(userUsecase, refreshTokenUsecase, jwtService)
	r.POST("/refreshtoken", refreshTokenController.RefreshToken)
}
