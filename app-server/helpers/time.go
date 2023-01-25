package helpers

import "time"

func Now() *time.Time {
	t := time.Now().UTC()
	return &t
}
