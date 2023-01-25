package tasks

import (
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type TaskListRes struct {
	Count int64         `json:"count"`
	Data  []models.Task `json:"data"`
}

func NewTaskList(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)

		tx := xctx.Db.Model(&models.Task{}).Where("user_id = ?", user.Id)
		tx = models.WithTaskPaging(tx, ctx.Query("_page"), ctx.Query("_limit"))
		tx = models.WithTaskCreationRange(tx, ctx.Query("created_at_start"), ctx.Query("created_at_end"))
		tx = models.WithTaskStatusQuery(tx, ctx.Query("status"))
		tx = models.WithTaskId(tx, ctx.Query("id"))

		var count int64
		if tx.Count(&count); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		var data []models.Task
		if tx.Order("created_at DESC").Find(&data); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		for i, task := range data {
			task.WithStatus()
			data[i] = task
		}

		res := TaskListRes{Count: count, Data: data}
		ctx.JSON(http.StatusOK, res)
	}
}
