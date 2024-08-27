package routers

import (
	"loan/domain"

	"github.com/gin-gonic/gin"
)

type ProfileRouter struct {
	profileController domain.ProfileHandler
	engine            *gin.Engine
}

func NewProfileRouter(p domain.ProfileHandler, engine *gin.Engine) domain.ProfileRouter {
	return &ProfileRouter{
		profileController: p,
		engine:            engine,
	}
}

func (p *ProfileRouter) InitProfileRoutes(auth *gin.RouterGroup) {
	// User profile routes

	auth.GET("users/profile/:user_id", p.profileController.FindProfile)
	auth.PUT("users/profile/", p.profileController.UpdateProfile)
	auth.POST("users/profile", p.profileController.SaveProfile)
}
