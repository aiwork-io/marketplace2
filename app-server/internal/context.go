package internal

import (
	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/models"
	"gorm.io/gorm"
)

type XContext struct {
	Configs *Configs
	Db      *gorm.DB
	Storage Storage
	Mailer  Mailer
}

func (ctx *XContext) Init() {
	ctx.Db.AutoMigrate(&models.User{}, &models.Task{}, &models.TaskAsset{})
}

func NewXContext() *XContext {
	configs := NewConfigs(NewConfigProvider(helpers.GetRootDir(), ".", "./secrets"))
	db := NewDbClient(configs.DbUri)
	storage := NewStorage(configs)
	mailer := NewMailer(configs)

	return &XContext{Configs: configs, Db: db, Storage: storage, Mailer: mailer}
}
