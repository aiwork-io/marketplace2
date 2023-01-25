package client_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"aiwork.io/marketplace/routes/client"
	"github.com/stretchr/testify/assert"
)

func TestNextTask(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	if err := helpers.TruncateDb(xctx.Db, &models.Task{}); err != nil {
		panic(err)
	}

	// seed user
	employer := internal.GenUser(xctx, internal.TEST_EMPLOYER_WALLET)
	employee := internal.GenUser(xctx, "")

	// ignore completed task
	before1h := time.Now().Add(-60 * time.Minute)
	internal.GenTask(xctx, "completed", employer, employee, &before1h)
	// ignore completed task of another user
	internal.GenTask(xctx, "completed", employer, internal.GenUser(xctx, ""), &before1h)
	// ignore processing task of another user
	internal.GenTask(xctx, "processing", employer, internal.GenUser(xctx, ""), &before1h)
	// pick this task
	before30m := time.Now().Add(-30 * time.Minute)
	willbeusetask := internal.GenTask(xctx, "pending", employer, employee, &before30m)
	// ignore newest task
	internal.GenTask(xctx, "pending", employer, employee, nil)

	server := internal.NewHttpServer()
	route := server.Group("/client")
	client.NewRouter(internal.NewFakeAuthMiddleware(route, employee), xctx)

	w := internal.DoTestReq(server, http.MethodGet, "/client/tasks/next", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var res models.Task
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, willbeusetask.Id, res.Id)
	assert.Equal(t, len(res.Assets), len(willbeusetask.Assets))
}
