package tasks

import (
	"log"
	"net/http"
	"os"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

func NewTaskResultsDownloadUrl(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		task := ctx.Keys["_task"].(*models.Task)
		fileKey := "results/" + task.Id + ".zip"

		// already exist in S3
		if exists := xctx.Storage.FileExist(fileKey); exists {
			url, err := xctx.Storage.PresignedUrl("GET", fileKey, 1)
			if err != nil {
				log.Println(err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
				return
			}
			if url != "" {
				ctx.JSON(http.StatusOK, map[string]interface{}{"url": url})
				return
			}
		}

		// start zip
		folderpath, err := task.ToTempFiles("/tmp")
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		target := folderpath + ".zip"
		if err := helpers.Zip(folderpath, target); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		if err := xctx.Storage.FileUpload(target, fileKey); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		// cleanup
		_ = os.RemoveAll(folderpath)
		_ = os.RemoveAll(target)

		url, err := xctx.Storage.PresignedUrl("GET", fileKey, 1)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}
		ctx.JSON(http.StatusOK, map[string]interface{}{"url": url})
	}
}
