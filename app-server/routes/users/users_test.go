package users_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/routes/users"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestProfileUpdate(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	user := internal.GenUser(xctx, "")
	updated := users.ProfileUpdateReq{
		Name:   gofakeit.Name(),
		Wallet: gofakeit.LetterN(64),
	}

	server := internal.NewHttpServer()
	route := server.Group("/users")
	users.NewRouter(internal.NewFakeAuthMiddleware(route, user), xctx)

	w := internal.DoTestReq(server, http.MethodPut, "/users/profile", helpers.MustMarshal(updated))

	assert.Equal(t, http.StatusOK, w.Code)

	var res users.ProfileUpdateRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, user.Id, res.Id)

	updateduser, _ := internal.GetUser(xctx, user.Id)
	assert.Equal(t, updateduser.Name, updated.Name)
	assert.Equal(t, updateduser.Wallet, updated.Wallet)
}

func TestPasswordUpdate(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	user := internal.GenUser(xctx, "")
	updated := users.PasswordUpdateReq{
		Password: gofakeit.LetterN(6),
	}
	// Make sure we didn't set the same password
	assert.NotNil(t, user.CheckPassword(updated.Password))

	server := internal.NewHttpServer()
	route := server.Group("/users")
	users.NewRouter(internal.NewFakeAuthMiddleware(route, user), xctx)

	w := internal.DoTestReq(server, http.MethodPut, "/users/password", helpers.MustMarshal(updated))

	assert.Equal(t, http.StatusOK, w.Code)

	var res users.PasswordUpdateRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, user.Id, res.Id)

	updateduser, _ := internal.GetUser(xctx, user.Id)
	assert.Nil(t, updateduser.CheckPassword(updated.Password))
}
