package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type VerificationReq struct {
	State   string                 `json:"state"  binding:"required"`
	Payload map[string]interface{} `json:"payload"` // Rest of the fields should go here.

}

type VerificationRes struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
}

func NewVerification(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req VerificationReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		payload, err := helpers.DecryptAES(xctx.Configs.Secret, req.State)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		var state models.AuthState
		if err := json.Unmarshal([]byte(payload), &state); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		if !state.Valid() {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "State was expired"})
			return
		}

		user, err := internal.GetUser(xctx, state.UserId)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not valid state"})
			return
		}

		switch state.Scenario {
		case models.AUTH_SCENARIO_ACCOUNT_CONFIRMATION:
			now := time.Now().UTC()
			user.VerifiedAt = &now
		case models.AUTH_SCENARIO_PASSWORD_RESET:
			password, _ := req.Payload["password"].(string)
			if password == "" {
				log.Printf("invalid payload %s", helpers.MustMarshal(req))
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Not valid payload"})
				return
			}
			user.SetPassword(password)
		default:
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not valid scenario"})
			return
		}

		if tx := xctx.Db.Save(user); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		token, err := user.SignToken(xctx.Configs.Secret)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
			return
		}

		res := VerificationRes{Id: user.Id, AccessToken: token}
		ctx.JSON(http.StatusOK, res)
	}
}
