package users

import (
	"net/http"
	"time"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

type GetProfileRes struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Wallet     string     `json:"wallet"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	VerifiedAt *time.Time `json:"verified_at"`
}

func NewGetProfile(xctx *internal.XContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)

		res := GetProfileRes{
			Id:         user.Id,
			Name:       user.Name,
			Email:      user.Email,
			Wallet:     user.Wallet,
			VerifiedAt: user.VerifiedAt,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		}
		ctx.JSON(http.StatusOK, res)
	}
}
