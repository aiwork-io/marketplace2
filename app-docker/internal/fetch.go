package internal

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

func NewFetch() *resty.Client {
	return resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= http.StatusInternalServerError
			},
		).
		SetTimeout(5 * time.Minute)
}
