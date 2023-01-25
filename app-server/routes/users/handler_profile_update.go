package users

import (
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type ProfileUpdateReq struct {
	Name   string `json:"name"  binding:"required"`
	Wallet string `json:"wallet"  binding:"required"`
}

type ProfileUpdateRes struct {
	Id string `json:"id"`
}

func NewProfileUpdate(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ProfileUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		user.Name = req.Name
		user.Wallet = req.Wallet

		if tx := xctx.Db.Save(user); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		res := ProfileUpdateRes{Id: user.Id}
		ctx.JSON(http.StatusOK, res)
	}
}
