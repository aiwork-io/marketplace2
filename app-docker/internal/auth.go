package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type AuthUser struct {
	Id          string `json:"id"`
	Role        int    `json:"role"`
	Wallet      string `json:"wallet"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type Authorize func() (*AuthUser, error)

func NewAuth(configs *Configs) Authorize {
	client := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= http.StatusInternalServerError
			},
		)
	return func() (*AuthUser, error) {
		req := client.R().SetBody(map[string]interface{}{
			"email":    configs.AuthUsername,
			"password": configs.AuthPassword,
		})

		uri := configs.ApiEndpoint + "/auth/login"
		res, err := req.Post(uri)
		if err != nil {
			return nil, err
		}

		ok := res.StatusCode() > 100 && res.StatusCode() < 300
		if !ok {
			return nil, fmt.Errorf("%s | login was failed", configs.AuthUsername)
		}

		var user AuthUser
		err = json.Unmarshal(res.Body(), &user)

		return &user, err
	}
}
