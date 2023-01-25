package helpers

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func GetRootDir() string {
	_, file, _, _ := runtime.Caller(1)
	curdir := path.Dir(file)
	dirpath, _ := filepath.Abs(curdir + "/..")
	return dirpath
}

func Lock(key string) (bool, error) {
	keypath := "/tmp/" + key
	_, err := os.Stat(keypath)

	// exists
	if err == nil {
		return false, nil
	}

	// not exists
	if errors.Is(err, os.ErrNotExist) {
		writeerr := os.WriteFile(keypath, []byte(keypath), 0644)
		return writeerr == nil, writeerr
	}

	// other error
	return false, err
}

func Unlock(key string) error {
	keypath := "/tmp/" + key
	return os.Remove(keypath)
}
