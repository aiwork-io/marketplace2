package internal

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDbClient(uri string) *gorm.DB {
	dialector := postgres.Open(uri)

	lconfig := logger.Config{
		SlowThreshold:             5 * time.Second, // Slow SQL threshold
		LogLevel:                  logger.Silent,   // Log level
		IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
		Colorful:                  false,           // Disable color
	}
	gormlogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		lconfig,
	)

	client, err := gorm.Open(dialector, &gorm.Config{Logger: gormlogger})
	if err != nil {
		panic(err)
	}

	return client
}
