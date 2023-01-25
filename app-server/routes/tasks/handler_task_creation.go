package tasks

import (
	"log"
	"net/http"
	"time"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type TaskCreationReq struct {
	Name       string `json:"name" binding:"required"`
	Category   string `json:"category" binding:"required"`
	PaymentTxn string `json:"payment_txn"`
}

type TaskCreationRes struct {
	Id        string     `json:"id"`
	UserId    string     `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

func NewTaskCreation(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req TaskCreationReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		task := models.Task{
			Id:          helpers.NewId("task"),
			UserId:      user.Id,
			Name:        req.Name,
			Category:    req.Category,
			PaymentTxn:  req.PaymentTxn,
			PaymentRate: xctx.Configs.TxnRate,
		}
		if task.PaymentTxn == "" {
			task.PaymentTxn = helpers.NewId("faketxn")
		}
		if tx := xctx.Db.Save(&task); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		res := TaskCreationRes{Id: task.Id, UserId: task.UserId, CreatedAt: task.CreatedAt}
		ctx.JSON(http.StatusOK, res)
	}
}
