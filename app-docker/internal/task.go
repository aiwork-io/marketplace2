package internal

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type taskstatus int64

const TASK_STATUS_PENDING taskstatus = 0
const TASK_STATUS_PROCESSING taskstatus = 100
const TASK_STATUS_COMPLETED taskstatus = 200

type Task struct {
	Id                string     `json:"id"`
	UserId            string     `json:"user_id"`
	Name              string     `json:"name"`
	Category          string     `json:"category"`
	PaymentTxn        string     `json:"payment_txn"`
	PaymentVerifiedAt *time.Time `json:"payment_verified_at"`
	PaymentProof      string     `json:"payment_proof"`
	CreatedAt         *time.Time `json:"created_at" `
	UpdatedAt         *time.Time `json:"updated_at" `

	ProcessingBy string     `json:"processing_by"`
	ProcessingAt *time.Time `json:"processing_at"`
	CompletedBy  string     `json:"completed_by"`
	CompletedAt  *time.Time `json:"completed_at"`

	RewardTxn  string     `json:"reward_txn"`
	RewardedAt *time.Time `json:"rewarded_at"`

	Assets []TaskAsset `json:"assets"` // ignore this field when write and read with struct
	Status taskstatus  `json:"status"` // ignore this field when write and read with struct
}

type TaskAsset struct {
	Id         string     `json:"id"`
	TaskId     string     `json:"task_id"`
	UserId     string     `json:"user_id"`
	FileBucket string     `json:"file_bucket"`
	FileKey    string     `json:"file_key"`
	Results    string     `json:"results"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`

	FileUrl string `json:"file_url"`
}

type TaskCompleteReq struct {
	AssetIds []string `json:"asset_ids" binding:"required"`
}

func NewTaskClient(configs *Configs) *TaskClient {
	return &TaskClient{
		auth:    NewAuth(configs),
		configs: configs,
		client:  NewFetch(),
	}
}

type TaskClient struct {
	auth    Authorize
	configs *Configs
	client  *resty.Client
}

func (tc *TaskClient) Get() (*Task, error) {
	user, err := tc.auth()
	if err != nil {
		return nil, err
	}

	req := tc.client.R().
		SetAuthToken(user.AccessToken)
	uri := tc.configs.ApiEndpoint + "/client/tasks/next"
	res, err := req.Get(uri)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, errors.New(res.Status() + ": " + string(res.Body()))
	}

	var task Task
	err = json.Unmarshal(res.Body(), &task)

	return &task, err
}

func (tc *TaskClient) Complete(task *Task, ids []string) (*Task, error) {
	user, err := tc.auth()
	if err != nil {
		return nil, err
	}

	req := tc.client.R().
		SetAuthToken(user.AccessToken).
		SetBody(TaskCompleteReq{AssetIds: ids})
	uri := tc.configs.ApiEndpoint + "/client/tasks/" + task.Id + "/submission"
	res, err := req.Put(uri)
	if err != nil {
		return nil, err
	}

	var updatedtask Task
	err = json.Unmarshal(res.Body(), &updatedtask)
	log.Printf("complete task %s with response %s", task.Id, string(res.Body()))

	return &updatedtask, err
}
