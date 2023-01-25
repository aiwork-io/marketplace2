package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/segmentio/ksuid"
)

func NewId(prefix string) string {
	return strings.Join([]string{prefix, ksuid.New().String()}, "_")
}

func NewRandomKey(length int) (string, error) {
	chars := make([]byte, 64)
	if _, err := rand.Read(chars); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(chars), nil
}

func MustMarshal(obj any) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func NewStorageKey(segments ...string) string {
	return strings.Join(segments, "/")
}
