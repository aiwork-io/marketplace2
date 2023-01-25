package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

func NewAuthJWTMiddleware(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// already authenticated
		if _, ok := ctx.Keys[models.CTX_KEY_USER].(*models.User); ok {
			ctx.Next()
			return
		}

		token := getAccessToken(ctx)
		if token == "" {
			log.Println(errors.New("no token"))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		u, err := models.ValidateToken(xctx.Configs.Secret, token)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := internal.GetUser(xctx, u.Id)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Keys[models.CTX_KEY_USER] = user
		ctx.Next()
	}
}

func getAccessToken(ctx *gin.Context) string {
	if token := ctx.Query("_access_token"); token != "" {
		return token
	}

	bearer := ctx.Request.Header.Get("authorization")
	segments := strings.Split(bearer, " ")

	if len(segments) == 2 {
		return segments[1]
	}

	return segments[0]
}
