package users

import (
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/routes/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.RouterGroup, xctx *internal.XContext) {
	route.Use(auth.NewAuthJWTMiddleware(xctx))
	route.GET("me", NewGetProfile(xctx))
	route.PUT("profile", NewProfileUpdate(xctx))
	route.PUT("password", NewPasswordUpdate(xctx))
}
