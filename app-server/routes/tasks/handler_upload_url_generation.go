package tasks

import (
	"errors"
	"log"
	"net/http"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type UploadUrlGenerationRes struct {
	Asset *models.TaskAsset `json:"asset"`
	Url   string            `json:"url"`
}

func NewUploadUrlGeneration(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filename := ctx.Query("filename")
		if filename == "" {
			log.Println(errors.New("no filename"))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
			return
		}

		task := ctx.Keys["_task"].(*models.Task)
		filekey := helpers.NewStorageKey(task.Id, filename)
		url, err := xctx.Storage.PresignedUrl("PUT", filekey, 1)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not generate upload url"})
			return
		}

		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		asset := models.TaskAsset{
			Id:         helpers.NewId("asset"),
			UserId:     user.Id,
			TaskId:     task.Id,
			FileBucket: xctx.Storage.CurrentBucket(),
			FileKey:    filekey,
		}
		if tx := xctx.Db.Save(&asset); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		ctx.JSON(http.StatusOK, UploadUrlGenerationRes{Asset: &asset, Url: url})
	}
}
