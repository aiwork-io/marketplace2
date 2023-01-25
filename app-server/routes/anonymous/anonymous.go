package anonymous

import (
	"net/http"

	"aiwork.io/marketplace/internal"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.RouterGroup, xctx *internal.XContext) {
	route.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"version": xctx.Configs.Version})
	})
}
