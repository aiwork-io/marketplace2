package models

import (
	"math"
	"os"
	"time"

	"gorm.io/gorm"
)

type taskstatus int64

const TASK_STATUS_PENDING taskstatus = 0
const TASK_STATUS_PROCESSING taskstatus = 100
const TASK_STATUS_COMPLETED taskstatus = 200

type Task struct {
	Id                string     `json:"id" gorm:"primaryKey;size:66"`
	UserId            string     `json:"user_id" gorm:"size:66"`
	Name              string     `json:"name" gorm:"size:512"`
	Category          string     `json:"category" gorm:"size:256"`
	PaymentRate       int        `json:"payment_rate"`
	PaymentTxn        string     `json:"payment_txn" gorm:"unique;size:128"`
	PaymentVerifiedAt *time.Time `json:"payment_verified_at"`
	CreatedAt         *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	ProcessingBy string     `json:"processing_by" gorm:"size:66"`
	ProcessingAt *time.Time `json:"processing_at"`
	CompletedBy  string     `json:"completed_by" gorm:"size:66"`
	CompletedAt  *time.Time `json:"completed_at"`

	RewardTxn  string     `json:"reward_txn" gorm:"size:66"`
	RewardedAt *time.Time `json:"rewarded_at"`

	Assets  []TaskAsset `json:"assets" gorm:"-"`  // ignore this field when write and read with struct
	Status  taskstatus  `json:"status" gorm:"-"`  // ignore this field when write and read with struct
	Credits float64     `json:"credits" gorm:"-"` // ignore this field when write and read with struct
}

func (t *Task) WithStatus() taskstatus {
	// Pending
	t.Status = TASK_STATUS_PENDING
	// override with processing
	if t.ProcessingAt != nil {
		t.Status = TASK_STATUS_PROCESSING
	}
	// override with completed
	if t.CompletedAt != nil {
		t.Status = TASK_STATUS_COMPLETED
	}
	return t.Status
}

func (t *Task) WithAssets(db *gorm.DB) error {
	tx := db.
		Model(&TaskAsset{}).
		Where("user_id = ? AND task_id = ?", t.UserId, t.Id).
		Find(&t.Assets)
	return tx.Error
}

func (t *Task) WithCredit() {
	if len(t.Assets) == 0 {
		return
	}

	credits := float64(len(t.Assets)) / float64(t.PaymentRate)
	t.Credits = math.Round(credits*100000) / 100000
}

func (t *Task) ToTempFiles(folder string) (string, error) {
	folderpath := folder + "/" + t.Id
	if err := os.MkdirAll(folderpath, os.ModePerm); err != nil {
		return folderpath, err
	}

	for _, asset := range t.Assets {
		filepath := folderpath + "/" + asset.Id + ".json"
		if err := os.WriteFile(filepath, []byte(asset.Results), 0644); err != nil {
			return folderpath, err
		}
	}

	return folderpath, nil
}

type TaskAsset struct {
	Id         string     `json:"id" gorm:"primaryKey;size:66"`
	TaskId     string     `json:"task_id" gorm:"size:66"`
	UserId     string     `json:"user_id" gorm:"size:66"`
	FileBucket string     `json:"file_bucket"`
	FileKey    string     `json:"file_key"`
	Results    string     `json:"results"`
	CreatedAt  *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	FileUrl string `json:"file_url" gorm:"-"`
}
