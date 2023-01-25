package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

func NewTxnValidator(endpoint, wallet, paymenttoken string) func(from, txn string, amount float64, source string) error {
	client := NewHttpClient()

	return func(from, txn string, amount float64, source string) error {
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		u.Path = "/check"

		query := u.Query()
		query.Add("from_wallet_id", from)
		query.Add("wallet_id", wallet)
		query.Add("transaction_id", txn)
		query.Add("source", source)
		query.Add("amount", fmt.Sprintf("%f", amount))
		query.Add("payment_token", paymenttoken)
		u.RawQuery = query.Encode()

		res, err := client.R().Get(u.String())
		if err != nil {
			log.Println(err)
			return err
		}

		var data map[string]bool
		if err := json.Unmarshal(res.Body(), &data); err != nil {
			log.Printf("could not parse data %s when submit to %s", string(res.Body()), u.String())
			return err
		}
		if !data["ok"] {
			log.Println("payment amount was not matched")
			return errors.New("payment amount was not matched")
		}

		return nil
	}
}

func NewTxnTransfer(endpoint, paymenttoken string) func(to string, amount float64) (string, error) {
	client := NewHttpClient()

	return func(to string, amount float64) (string, error) {
		u, err := url.Parse(endpoint)
		if err != nil {
			return "", err
		}
		u.Path = "/transfer"

		body := map[string]interface{}{
			"recipient_wallet": to,
			"amount":           amount,
		}
		res, err := client.R().SetBody(body).Post(u.String())
		if err != nil {
			log.Println(err)
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(res.Body(), &data); err != nil {
			log.Printf("could not parse data %s when submit to %s", string(res.Body()), u.String())
			return "", err
		}

		txnhash, ok := data["transactionHash"].(string)
		if !ok || txnhash == "" {
			log.Printf("transfer was failed with respone %s", string(res.Body()))
			return "", errors.New("transfer was failed")
		}

		return txnhash, nil
	}
}

func NewHttpClient() *resty.Client {
	return resty.New().
		SetTimeout(time.Duration(3) * time.Second).
		SetRetryCount(1).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= http.StatusInternalServerError
			},
		).
		SetTimeout(5 * time.Minute)
}
