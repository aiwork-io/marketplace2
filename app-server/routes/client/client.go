package client

import (
	"fmt"
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"aiwork.io/marketplace/routes/auth"
	"aiwork.io/marketplace/routes/tasks"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.RouterGroup, xctx *internal.XContext) {
	route.Use(auth.NewAuthJWTMiddleware(xctx))
	route.GET("tasks", NewTaskList(xctx))
	route.POST(
		"tasks/:task_id/assets/:asset_id/results",
		tasks.NewGetTaskId(xctx),
		NewTaskAuthority(xctx),
		NewTaskAssetSubmission(xctx),
	)
	route.PUT(
		"tasks/:task_id/submission",
		tasks.NewGetTaskId(xctx),
		NewTaskAuthority(xctx),
		NewTaskComplete(xctx),
	)
	route.GET("tasks/next", NewNextTask(xctx))
}

func NewTaskAuthority(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		task := ctx.Keys["_task"].(*models.Task)

		if task.ProcessingBy != user.Id {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not your task"})
			return
		}

		if task.Status != models.TASK_STATUS_PROCESSING {
			err := fmt.Errorf("task status is not valid | expected %d, acutal: %d", models.TASK_STATUS_PROCESSING, task.Status)
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.Next()
	}
}
