package internal

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
)

const TEST_TXN_ID string = "0xa9e22075696e614d13cb4e3fac621e63e3ef2162496a1b7e287f99f6bb142b4b"
const TEST_EMPLOYER_WALLET = "0xA7E9eDbD4d311613963764BD65870c863967A92a"

func DoTestReq(r *gin.Engine, method, url string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	httpreq, _ := http.NewRequest(method, url, bytes.NewReader(body))
	r.ServeHTTP(w, httpreq)
	return w
}

func GenTask(
	xctx *XContext,
	status string,
	employer *models.User,
	employee *models.User,
	created *time.Time,
) *models.Task {
	now := time.Now().UTC()
	task := models.Task{
		Id:       helpers.NewId("task"),
		UserId:   employer.Id,
		Name:     gofakeit.Animal(),
		Category: gofakeit.EmojiCategory(),
		Assets:   []models.TaskAsset{},
	}
	if created != nil {
		now = *created
		task.CreatedAt = created
	}

	switch status {
	case "processing":
		if employee != nil {
			task.ProcessingBy = employee.Id
		}
		task.ProcessingAt = &now
	case "completed":
		if employee != nil {
			task.ProcessingBy = employee.Id
			task.CompletedBy = employee.Id
		}
		task.CompletedAt = &now
	default:
		task.PaymentVerifiedAt = &now
	}

	if tx := xctx.Db.Create(&task); tx.Error != nil {
		panic(tx.Error)
	}

	count := gofakeit.IntRange(3, 9)
	for i := 0; i < count; i++ {
		task.Assets = append(task.Assets, models.TaskAsset{
			Id:         helpers.NewId("asset"),
			UserId:     employer.Id,
			TaskId:     task.Id,
			FileBucket: gofakeit.LetterN(9),
			FileKey:    gofakeit.LetterN(3) + gofakeit.FileExtension(),
		})
	}
	if tx := xctx.Db.Create(&task.Assets); tx.Error != nil {
		panic(tx.Error)
	}

	return &task
}

func GenUser(xctx *XContext, wallet string) *models.User {
	user := models.User{
		Id:     helpers.NewId("user"),
		Name:   gofakeit.Name(),
		Wallet: gofakeit.LetterN(64),
		Email:  gofakeit.Email(),
	}
	if wallet != "" {
		user.Wallet = wallet
	}
	user.SetPassword(gofakeit.Password(true, true, true, false, false, 6))
	if tx := xctx.Db.Create(&user); tx.Error != nil {
		panic(tx.Error)
	}
	return &user
}

func GenRealworldUser(xctx *XContext) *models.User {
	if user, err := GetUser(xctx, xctx.Configs.TestEmail); err == nil {
		return user
	}

	user := models.User{
		Id:     helpers.NewId("user"),
		Name:   gofakeit.Name(),
		Wallet: TEST_EMPLOYER_WALLET,
		Email:  xctx.Configs.TestEmail,
	}
	user.SetPassword(gofakeit.Password(true, true, true, false, false, 6))
	if tx := xctx.Db.Create(&user); tx.Error != nil {
		panic(tx.Error)
	}
	return &user
}

func GetUser(xctx *XContext, id string) (*models.User, error) {
	user := models.User{}
	if tx := xctx.Db.Where("id = ? OR email = ?", id, id).First(&user); tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func NewFakeAuthMiddleware(route *gin.RouterGroup, user *models.User) *gin.RouterGroup {
	route.Use(func(ctx *gin.Context) {
		ctx.Keys[models.CTX_KEY_USER] = user
	})

	return route
}
