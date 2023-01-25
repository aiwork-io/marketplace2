package client

import (
	"errors"
	"io"
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type EngineTask struct {
	Source  string                 `json:"source"`
	Context *EngineTaskContenxt    `json:"context"`
	Next    []string               `json:"next"`
	Prev    []string               `json:"prev"`
	Data    map[string]interface{} `json:"data"`
	Action  string                 `json:"action"`
}

type EngineTaskContenxt struct {
	ProjectId      string   `json:"project_id"`
	Timecode       string   `json:"timecode"`
	Debug          bool     `json:"debug"`
	InterestRegion []string `json:"interest_region"`
	ObjectFilter   []string `json:"object_filter"`
}

func NewTaskAssetSubmission(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		task := ctx.Keys["_task"].(*models.Task)

		assetId := ctx.Param("asset_id")
		if assetId == "" {
			log.Println(errors.New("no asset"))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No asset id"})
			return
		}

		asset := models.TaskAsset{Id: assetId, TaskId: task.Id}
		if tx := xctx.Db.First(&asset); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No asset was found"})
			return
		}

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs. Please check your inputs"})
			return
		}
		asset.Results = string(body)
		if tx := xctx.Db.Save(&asset); tx.Error != nil {
			log.Println(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		ctx.JSON(http.StatusOK, map[string]string{"user_id": user.Id, "task_id": task.Id})
	}
}
