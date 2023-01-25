package auth

import (
	"fmt"
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type PasswordResetReq struct {
	Email string `json:"email"  binding:"required"`
	Code  int64  `json:"code"`
}

type PasswordResetRes struct {
	Code int64 `json:"code"`
}

func NewPasswordReset(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PasswordResetReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user, err := internal.GetUser(xctx, req.Email)
		if err != nil {
			log.Println(fmt.Errorf("no user was found with email %s", req.Email))
			ctx.JSON(http.StatusOK, PasswordResetRes{})
			return
		}

		if err := SendResetPassword(xctx, user); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, PasswordResetRes{})
			return
		}

		ctx.JSON(http.StatusOK, PasswordResetRes{Code: req.Code})
	}
}

func SendResetPassword(xctx *internal.XContext, user *models.User) error {
	uri, err := GenStateUrl(user, xctx.Configs.AuthVerifictionResetEndpoint, xctx.Configs.Secret, models.AUTH_SCENARIO_PASSWORD_RESET)
	if err != nil {
		return err
	}
	to := user.Email
	subject := "Password Recovery Confirmation"
	body := "Please follow this url to reset your password: " + uri

	return xctx.Mailer.Send(to, subject, body)
}
