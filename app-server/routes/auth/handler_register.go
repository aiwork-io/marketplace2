package auth

import (
	"log"
	"net/http"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Name     string `json:"name"  binding:"required"`
	Wallet   string `json:"wallet"  binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}
type RegisterRes struct {
	Id          string `json:"id"`
	Wallet      string `json:"wallet"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func NewRegister(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RegisterReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user := models.User{
			Id:     helpers.NewId("user"),
			Name:   req.Name,
			Email:  req.Email,
			Wallet: req.Wallet,
		}
		user.SetPassword(req.Password)
		if tx := xctx.Db.Save(&user); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		if err := SendAccountConfirmation(xctx, &user); err != nil {
			log.Println(err)
		}

		token, err := user.SignToken(xctx.Configs.Secret)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
			return
		}

		res := RegisterRes{
			Id:          user.Id,
			Email:       user.Email,
			Wallet:      user.Wallet,
			AccessToken: token,
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func SendAccountConfirmation(xctx *internal.XContext, user *models.User) error {
	enable := xctx.Configs.AuthVerifictionAccountEnable == 1
	if !enable {
		return nil
	}

	uri, err := GenStateUrl(user, xctx.Configs.AuthVerifictionAccountEndpoint, xctx.Configs.Secret, models.AUTH_SCENARIO_ACCOUNT_CONFIRMATION)
	if err != nil {
		return err
	}
	to := user.Email
	subject := "New Account Confirmation"
	body := "Please follow this url to confirm your account: " + uri

	return xctx.Mailer.Send(to, subject, body)
}
