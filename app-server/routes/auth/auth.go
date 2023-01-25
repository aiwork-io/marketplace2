package auth

import (
	"net/http"
	"net/url"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"github.com/gin-gonic/gin"
)

func NewRouter(route *gin.RouterGroup, xctx *internal.XContext) {
	route.POST("recovery/password", NewPasswordReset(xctx))
	route.POST("recovery/verification", NewVerification(xctx))

	route.POST("", NewRegister(xctx))
	route.POST("login", NewSignIn(xctx), func(ctx *gin.Context) {
		res := ctx.Keys["_signin_res"]
		ctx.JSON(http.StatusOK, res)
	})
}

func GenStateUrl(user *models.User, endpoint, secret, scenario string) (string, error) {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	payload := helpers.MustMarshal(models.NewAuthState(user.Id, scenario))
	state, err := helpers.EncryptAES(secret, string(payload))
	if err != nil {
		return "", err
	}

	query := uri.Query()
	query.Set("state", state)
	uri.RawQuery = query.Encode()

	return uri.String(), nil
}
