package admin

import (
	"fmt"
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

func NewAuthJWTMiddleware(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if !user.IsAdmin() {
			log.Println(fmt.Errorf("user %s does not have admin privilege", user.Id))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
