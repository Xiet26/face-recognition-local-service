package model

import (
	"time"
)

type AttendTemp struct {
	BatchID        string          `json:"batchID" bson:"batchID"`
	Time           time.Time       `json:"time" bson:"time"`
	StudentAttends []StudentAttend `json:"studentAttends" bson:"studentAttends"`
}

type StudentAttend struct {
	ImageFace string `json:"imageFace"`
	FaceID    int    `json:"faceID"`
}
