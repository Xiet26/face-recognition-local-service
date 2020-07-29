package model

import "gopkg.in/mgo.v2/bson"

type Face struct {
	ID        bson.ObjectId `json:"id"bson:"_id"`
	Vector    [128]float32  `json:"vector" bson:"vector"`
	FaceID    int32         `json:"faceID" bson:"faceID"`
}
