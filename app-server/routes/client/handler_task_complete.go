package client

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type TaskCompleteReq struct {
	AssetIds []string `json:"asset_ids" binding:"required"`
}

func NewTaskComplete(xctx *internal.XContext) gin.HandlerFunc {
	transfer := internal.NewTxnTransfer(xctx.Configs.TxnValidatorEndpoint, xctx.Configs.TxnValidatorPaymentToken)

	return func(ctx *gin.Context) {
		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		task := ctx.Keys["_task"].(*models.Task)

		var req TaskCompleteReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Printf("id:%s, err:%s", task.Id, err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}

		// BEGIN update
		taskAssetIds := lo.Map(task.Assets, func(asset models.TaskAsset, _ int) string { return asset.Id })
		sort.Strings(taskAssetIds)
		diffReq, diffTask := lo.Difference(req.AssetIds, taskAssetIds)
		if len(diffReq) > 0 || len(diffTask) > 0 {
			err := fmt.Errorf("id: %s missing asset | expected %d, acutal: %d", task.Id, len(taskAssetIds), len(req.AssetIds))
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatesql := "UPDATE tasks SET completed_by = ?, completed_at = CURRENT_TIMESTAMP WHERE processing_by = ? AND id = ?"
		if tx := xctx.Db.Exec(updatesql, user.Id, user.Id, task.Id); tx.Error != nil {
			log.Printf("id:%s, err:%s", task.Id, tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		updatedtask := models.Task{Id: task.Id}
		if tx := xctx.Db.First(&updatedtask); tx.Error != nil {
			log.Printf("id:%s, err:%s", task.Id, tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}
		// END update

		rewardAmount := task.Credits
		log.Printf("id:%s, to:%s, amount:%f", updatedtask.Id, user.Wallet, rewardAmount)
		// BEGIN reward
		if updatedtask.RewardTxn == "" {
			txn, err := transfer(user.Wallet, rewardAmount)
			log.Printf("id:%s, reward_txn:%s", updatedtask.Id, txn)

			if err != nil {
				log.Printf("id:%s, err:%s", updatedtask.Id, err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
				return
			}
			updatedtask.RewardTxn = txn
			updatedtask.RewardedAt = helpers.Now()
			if tx := xctx.Db.Save(&updatedtask); tx.Error != nil {
				log.Printf("id:%s, err:%s", updatedtask.Id, tx.Error.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
				return
			}
		}
		// END reward

		updatedtask.WithStatus()
		if err := updatedtask.WithAssets(xctx.Db); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}
		ctx.JSON(http.StatusOK, updatedtask)
	}
}
