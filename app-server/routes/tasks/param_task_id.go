package tasks

import (
	"errors"
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

func NewGetTaskId(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("task_id")
		if id == "" {
			log.Println(errors.New("no task id"))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
			return
		}

		task := models.Task{}
		tx := xctx.Db.Model(&models.Task{}).Where("id = ?", id).First(&task)
		if tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Task was not found"})
			return
		}

		if err := task.WithAssets(xctx.Db); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		task.WithStatus()
		task.WithCredit()
		ctx.Keys["_task"] = &task
		ctx.Next()
	}
}
