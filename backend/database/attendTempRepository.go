package database

import "xiet26/face-recognition-local-service/backend/model"

type AttendTempMongoRepository interface {
	Create(temp model.AttendTemp) error
	ReadByBatchID(batchID string) ([]model.AttendTemp, error)
	DeleteByBatchID(batchID string) error
}
