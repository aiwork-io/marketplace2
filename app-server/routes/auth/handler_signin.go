package auth

import (
	"log"
	"net/http"
	"time"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type SignInReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type SignInRes struct {
	Id          string     `json:"id"`
	Role        int        `json:"role"`
	Wallet      string     `json:"wallet"`
	Email       string     `json:"email"`
	VerifiedAt  *time.Time `json:"verified_at"`
	AccessToken string     `json:"access_token"`
}

func NewSignIn(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SignInReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user := models.User{}
		tx := xctx.Db.Model(&models.User{}).Where("email = ?", req.Email).First(&user)
		if tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		}

		if err := user.CheckPassword(req.Password); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		}

		token, err := user.SignToken(xctx.Configs.Secret)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not generate access token"})
			return
		}

		res := SignInRes{
			Id:          user.Id,
			Role:        user.Role,
			Email:       user.Email,
			Wallet:      user.Wallet,
			VerifiedAt:  user.VerifiedAt,
			AccessToken: token,
		}
		ctx.Keys["_signin_res"] = &res
		ctx.Next()
	}
}
