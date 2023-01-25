package tasks_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/models"
	"aiwork.io/marketplace/routes/tasks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreation(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	server := internal.NewHttpServer()
	route := server.Group("/task")
	user := internal.GenUser(xctx, internal.TEST_EMPLOYER_WALLET)
	tasks.NewRouter(internal.NewFakeAuthMiddleware(route, user), xctx)

	req := tasks.TaskCreationReq{
		Name:     gofakeit.Animal(),
		Category: gofakeit.EmojiCategory(),
	}

	w := internal.DoTestReq(server, http.MethodPost, "/task", helpers.MustMarshal(req))
	assert.Equal(t, http.StatusOK, w.Code)

	var res tasks.TaskCreationRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NotNil(t, res.Id)
	assert.NotNil(t, res.CreatedAt)
	assert.Equal(t, res.UserId, user.Id)
}

func TestTaskList(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	if err := helpers.TruncateDb(xctx.Db, &models.Task{}); err != nil {
		panic(err)
	}

	server := internal.NewHttpServer()
	route := server.Group("/task")
	user := internal.GenUser(xctx, internal.TEST_EMPLOYER_WALLET)
	tasks.NewRouter(internal.NewFakeAuthMiddleware(route, user), xctx)

	// seed user
	employer := internal.GenUser(xctx, internal.TEST_EMPLOYER_WALLET)
	employee := internal.GenUser(xctx, "")

	// before
	before1h30 := time.Now().Add(-90 * time.Minute)
	internal.GenTask(xctx, "processing", employer, employee, &before1h30)
	before1h := time.Now().Add(-60 * time.Minute)
	internal.GenTask(xctx, "processing", employer, employee, &before1h)
	before30m := time.Now().Add(-30 * time.Minute)
	internal.GenTask(xctx, "processing", employer, employee, &before30m)

	// now
	now := time.Now()
	internal.GenTask(xctx, "completed", employer, employee, &now)

	// after
	after1h30 := time.Now().Add(90 * time.Minute)
	internal.GenTask(xctx, "processing", employer, employee, &after1h30)
	after1h := time.Now().Add(60 * time.Minute)
	internal.GenTask(xctx, "processing", employer, employee, &after1h)
	after30m := time.Now().Add(30 * time.Minute)
	internal.GenTask(xctx, "processing", employer, employee, &after30m)

	// filters
	before45m := time.Now().Add(-45 * time.Minute).UTC().Format(time.RFC3339Nano)
	after45m := time.Now().Add(45 * time.Minute).UTC().Format(time.RFC3339Nano)

	uri := "/task?created_at_start=" + before45m + "&created_at_end=" + after45m + "&status=processing"
	w := internal.DoTestReq(server, http.MethodGet, uri, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var res tasks.TaskListRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, int64(2), res.Count)

	for _, task := range res.Data {
		assert.NotNil(t, task.ProcessingAt)
		assert.Nil(t, task.CompletedAt)
	}
}

func TestPaymentUpdate(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	// seed user
	employer := internal.GenUser(xctx, internal.TEST_EMPLOYER_WALLET)
	employee := internal.GenUser(xctx, "")

	server := internal.NewHttpServer()
	route := server.Group("/task")
	tasks.NewRouter(internal.NewFakeAuthMiddleware(route, employer), xctx)

	tr := internal.GenTask(xctx, "", employer, employee, nil)
	req := tasks.PaymentUpdateReq{
		PaymentTxn: internal.TEST_TXN_ID,
	}

	w := internal.DoTestReq(server, http.MethodPatch, "/task/"+tr.Id+"/payment", helpers.MustMarshal(req))
	assert.Equal(t, http.StatusOK, w.Code)

	var res tasks.PaymentUpdateRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NotNil(t, res.Id)
	assert.Equal(t, req.PaymentTxn, res.PaymentTxn)
	assert.NotNil(t, res.PaymentVerifiedAt)
}

func TestUploadUrlGeneration(t *testing.T) {
	xctx := internal.NewXContext()
	xctx.Init()

	// seed user
	employer := internal.GenUser(xctx, internal.TEST_EMPLOYER_WALLET)
	employee := internal.GenUser(xctx, "")

	server := internal.NewHttpServer()
	route := server.Group("/task")
	tasks.NewRouter(internal.NewFakeAuthMiddleware(route, employer), xctx)

	tr := internal.GenTask(xctx, "", employer, employee, nil)
	w := internal.DoTestReq(server, http.MethodGet, "/task/"+tr.Id+"/images/upload-url?filename=doggo.png", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var res tasks.UploadUrlGenerationRes
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.NotEmpty(t, res.Url)
	assert.NotNil(t, res.Asset)
	assert.NotEmpty(t, res.Asset.Id)
	assert.NotEmpty(t, res.Asset.FileBucket)
	assert.NotEmpty(t, res.Asset.FileKey)
}
