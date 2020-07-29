package database

import "xiet26/face-recognition-local-service/backend/model"

type FaceMongoRepository interface {
	Create(face model.Face) error
	ReadByFaceID(faceID int32) (model.Face, error)
	ReadByMultiFaceID(faceID []int32) ([]model.Face, error)
	DeleteByFaceID(faceID int32) error
	IsExist(faceID int32) bool
}
