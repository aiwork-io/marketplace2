package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"aiwork.io/marketplace/routes/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	server := internal.NewHttpServer()
	auth.NewRouter(server.Group("/auth"), xctx)
	w := internal.DoTestReq(server, http.MethodPost, "/auth", helpers.MustMarshal(newAccount()))

	assert.Equal(t, http.StatusOK, w.Code)

	var res auth.RegisterRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NotEmpty(t, res.Id)
	assert.NotEmpty(t, res.AccessToken)
}

func TestSignIn(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	server := internal.NewHttpServer()
	auth.NewRouter(server.Group("/auth"), xctx)

	account := newAccount()
	// register first
	internal.DoTestReq(server, http.MethodPost, "/auth", helpers.MustMarshal(account))

	// then sign in
	cred := auth.SignInReq{
		Email:    account.Email,
		Password: account.Password,
	}
	w := internal.DoTestReq(server, http.MethodPost, "/auth/login", helpers.MustMarshal(cred))
	var res auth.SignInRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.NotEmpty(t, res.Id)
	assert.NotEmpty(t, res.AccessToken)
}

func TestAuthMiddleware(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	server := internal.NewHttpServer()
	auth.NewRouter(server.Group("/auth"), xctx)
	server.GET("/private", auth.NewAuthJWTMiddleware(xctx), func(ctx *gin.Context) {
		user := ctx.Keys[models.CTX_KEY_USER].(*models.User)
		ctx.String(http.StatusOK, user.Id)
	})

	account := newAccount()
	registerw := internal.DoTestReq(server, http.MethodPost, "/auth", helpers.MustMarshal(account))
	var registerres auth.RegisterRes
	_ = json.Unmarshal(registerw.Body.Bytes(), &registerres)

	w := httptest.NewRecorder()
	httpreq, _ := http.NewRequest(http.MethodGet, "/private", nil)
	httpreq.Header.Set("Authorization", "Bearer "+registerres.AccessToken)
	server.ServeHTTP(w, httpreq)

	assert.Equal(t, registerres.Id, w.Body.String())

}

func newAccount() auth.RegisterReq {
	return auth.RegisterReq{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Wallet:   gofakeit.LetterN(66),
		Password: gofakeit.LetterN(8),
	}
}

func TestPasswordRecoveryReset(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	_ = internal.GenRealworldUser(xctx)

	server := internal.NewHttpServer()
	route := server.Group("/users")
	auth.NewRouter(route, xctx)

	req := auth.PasswordResetReq{
		Email: xctx.Configs.TestEmail,
		Code:  time.Now().Unix(),
	}
	w := internal.DoTestReq(server, http.MethodPost, "/auth/recovery/password", helpers.MustMarshal(req))

	assert.Equal(t, http.StatusOK, w.Code)

	var res auth.PasswordResetRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, req.Code, res.Code)

}
