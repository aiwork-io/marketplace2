package users

import (
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type PasswordUpdateReq struct {
	Password string `json:"password"  binding:"required"`
}

type PasswordUpdateRes struct {
	Id string `json:"id"`
}

func NewPasswordUpdate(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req PasswordUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		user.SetPassword(req.Password)

		if tx := xctx.Db.Save(user); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		res := PasswordUpdateRes{Id: user.Id}
		ctx.JSON(http.StatusOK, res)
	}
}
