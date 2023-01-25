package tasks

import (
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/routes/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.RouterGroup, xctx *internal.XContext) {
	route.Use(auth.NewAuthJWTMiddleware(xctx))
	route.POST("", NewTaskCreation(xctx))
	route.GET("", NewTaskList(xctx))
	route.GET(":task_id", NewGetTaskId(xctx), NewGetTask(xctx))
	route.GET(":task_id/results", NewGetTaskId(xctx), NewTaskResultsDownloadUrl(xctx))
	route.PATCH(":task_id/payment", NewGetTaskId(xctx), NewPaymentUpdate(xctx))
	route.GET(":task_id/images/upload-url", NewGetTaskId(xctx), NewUploadUrlGeneration(xctx))
}
