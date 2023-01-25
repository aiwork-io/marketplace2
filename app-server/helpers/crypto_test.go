package helpers_test

import (
	"testing"

	"aiwork.io/marketplace/helpers"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key := gofakeit.LetterN(32)
	value := gofakeit.UUID()

	encrypted, err := helpers.EncryptAES(key, value)
	assert.Nil(t, err)
	decrypted, err := helpers.DecryptAES(key, encrypted)
	assert.Nil(t, err)

	assert.Equal(t, value, decrypted)
}
