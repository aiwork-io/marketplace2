package internal

import (
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func NewEngineClient(configs *Configs) *EngineClient {
	return &EngineClient{
		auth:    NewAuth(configs),
		configs: configs,
		client:  NewFetch(),
	}
}

type EngineTask struct {
	Source  string                 `json:"source"`
	Context *EngineTaskContenxt    `json:"context"`
	Next    []string               `json:"next"`
	Prev    []string               `json:"prev"`
	Data    map[string]interface{} `json:"data"`
	Action  string                 `json:"action"`
}

type EngineTaskContenxt struct {
	ProjectId      string   `json:"project_id"`
	Timecode       string   `json:"timecode"`
	Debug          bool     `json:"debug"`
	InterestRegion []string `json:"interest_region"`
	ObjectFilter   []string `json:"object_filter"`
}

type EngineClient struct {
	auth    Authorize
	configs *Configs
	client  *resty.Client
}

func (ec *EngineClient) Push(asset *TaskAsset) error {
	user, err := ec.auth()
	if err != nil {
		return err
	}

	callback := ec.configs.ApiEndpoint + "/client/tasks/" + asset.TaskId + "/assets/" + asset.Id + "/results?_access_token=" + user.AccessToken

	body := EngineTask{
		Source: asset.FileUrl,
		Context: &EngineTaskContenxt{
			ProjectId:      asset.TaskId,
			Timecode:       "timecode",
			InterestRegion: []string{},
			ObjectFilter:   []string{},
			Debug:          true,
		},
		Data:   map[string]interface{}{},
		Action: "process",
		Next:   []string{"post:" + callback},
		Prev:   []string{},
	}

	req := ec.client.R().SetBody(body)
	uri := ec.configs.EngineEndpoint + "/receive"
	res, err := req.Post(uri)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New(res.Status() + ": " + string(res.Body()))
	}

	return nil
}
