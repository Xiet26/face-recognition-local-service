package database

import "xiet26/goface/face-management/model"

type AttendTempMongoRepository interface {
	Create(temp model.AttendTemp) error
	ReadByBatchID(batchID string) ([]model.AttendTemp, error)
	DeleteByBatchID(batchID string) error
}
