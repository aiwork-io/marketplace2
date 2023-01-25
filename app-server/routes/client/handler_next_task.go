package client

import (
	"log"
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewNextTask(xctx *internal.XContext) gin.HandlerFunc {
	takesql := `
UPDATE
	tasks
SET
	processing_by = ?,
	processing_at = CURRENT_TIMESTAMP
WHERE
	tasks.id in (
		SELECT
			t.id
		FROM
			tasks t
		WHERE
			t.payment_verified_at IS NOT NULL
			AND (t.completed_at IS NULL AND t.processing_at IS NULL)
			AND NOT (t.processing_by = ? AND t.processing_at IS NOT NULL)
		ORDER BY
			t.created_at
		LIMIT 1
	)
returning tasks.*`

	getsql := `
SELECT
	t.*
FROM
	tasks t
WHERE
	t.processing_by = ? and t.completed_at IS NULL
LIMIT 1`

	return func(ctx *gin.Context) {
		var task models.Task

		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		err := xctx.Db.Transaction(func(tx *gorm.DB) error {
			if subtx := tx.Raw(takesql, user.Id, user.Id).Scan(&task); subtx.Error != nil {
				return subtx.Error
			}

			// Find task that is processing by current user so the update sql didn't return data
			if subtx := tx.Raw(getsql, user.Id).Scan(&task); subtx.Error != nil {
				return subtx.Error
			}

			// return nil will commit the whole transaction
			return nil
		})

		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

		if task.Id == "" {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No task was found!"})
			return
		}

		if err := task.WithAssets(xctx.Db); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
			return
		}

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
