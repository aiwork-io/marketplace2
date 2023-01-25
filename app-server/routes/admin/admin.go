package admin

import (
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/routes/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.RouterGroup, xctx *internal.XContext) {
	route.POST("login", auth.NewSignIn(xctx), func(ctx *gin.Context) {
		res := ctx.Keys["_signin_res"]
		ctx.JSON(http.StatusOK, res)
	})
	route.POST("users", auth.NewAuthJWTMiddleware(xctx), NewAuthJWTMiddleware(xctx), auth.NewRegister(xctx))
}
