package models

import (
	"math"
	"strconv"

	"gorm.io/gorm"
)

func WithTaskPaging(tx *gorm.DB, p, l string) *gorm.DB {
	newtx := tx
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 20
	}
	newtx = newtx.Limit(limit)

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 0
	}
	offset := int(math.Max(float64((page-1)*limit), float64(0)))
	newtx = newtx.Offset(offset)
	return newtx
}

func WithTaskCreationRange(tx *gorm.DB, start, end string) *gorm.DB {
	newtx := tx

	if start != "" {
		newtx = newtx.Where("created_at > ?", start)
	}
	if end != "" {
		newtx = newtx.Where("created_at < ?", end)
	}

	return newtx
}

func WithTaskStatusQuery(tx *gorm.DB, status string) *gorm.DB {
	statusenum, err := strconv.Atoi(status)
	if err != nil {
		return tx
	}

	if statusenum == int(TASK_STATUS_PENDING) {
		return tx.Where("processing_at IS NULL and completed_at IS NULL")
	}
	if statusenum == int(TASK_STATUS_PROCESSING) {
		return tx.Where("processing_at IS NOT NULL and completed_at IS NULL")
	}
	if statusenum == int(TASK_STATUS_COMPLETED) {
		return tx.Where("completed_at IS NOT NULL")
	}
	return tx
}

func WithTaskId(tx *gorm.DB, id string) *gorm.DB {
	if id == "" {
		return tx
	}

	return tx.Where("id = ?", id)
}
