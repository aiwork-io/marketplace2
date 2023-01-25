package tasks

import (
	"log"
	"net/http"
	"time"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type PaymentUpdateReq struct {
	PaymentTxn string `json:"payment_txn" binding:"required"`
}

type PaymentUpdateRes struct {
	Id                string     `json:"id"`
	PaymentTxn        string     `json:"payment_txn"`
	PaymentVerifiedAt *time.Time `json:"payment_verified_at"`
}

func NewPaymentUpdate(xctx *internal.XContext) gin.HandlerFunc {
	check := internal.NewTxnValidator(
		xctx.Configs.TxnValidatorEndpoint,
		xctx.Configs.TxnValidatorDestAddress,
		xctx.Configs.TxnValidatorPaymentToken,
	)

	return func(ctx *gin.Context) {
		var req PaymentUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		task := ctx.Keys["_task"].(*models.Task)
		task.PaymentTxn = req.PaymentTxn

		if err := check(user.Wallet, task.PaymentTxn, task.Credits, task.Id); err == nil {
			now := time.Now().UTC()
			task.PaymentVerifiedAt = &now
		}

		if tx := xctx.Db.Save(&task); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		res := PaymentUpdateRes{
			Id:                task.Id,
			PaymentTxn:        task.PaymentTxn,
			PaymentVerifiedAt: task.PaymentVerifiedAt,
		}
		ctx.JSON(http.StatusOK, res)
	}
}
