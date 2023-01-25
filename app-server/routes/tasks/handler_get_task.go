package tasks

import (
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

func NewGetTask(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		task := ctx.Keys["_task"].(*models.Task)

		for i, asset := range task.Assets {
			url, err := xctx.Storage.PresignedUrl("GET", asset.FileKey, 24)
			if err != nil {
				log.Println(err)
				continue
			}
			asset.FileUrl = url
			task.Assets[i] = asset
		}

		ctx.JSON(http.StatusOK, task)
	}
}
