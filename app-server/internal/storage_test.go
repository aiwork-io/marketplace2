package internal_test

import (
	"testing"

	"aiwork.io/marketplace/helpers"
	"aiwork.io/marketplace/internal"
	"github.com/stretchr/testify/assert"
)

func TestGetPresignedUrl(t *testing.T) {
	filekey := helpers.NewStorageKey(helpers.NewId("task"), "doggo.png")

	configs := internal.NewConfigs(internal.NewConfigProvider(".", "./secrets"))
	storage := internal.NewStorage(configs)

	puturl, err := storage.PresignedUrl("PUT", filekey, 1)
	assert.Empty(t, err)
	assert.NotEmpty(t, puturl)

	geturl, err := storage.PresignedUrl("GET", filekey, 1)
	assert.Empty(t, err)
	assert.NotEmpty(t, geturl)
}
