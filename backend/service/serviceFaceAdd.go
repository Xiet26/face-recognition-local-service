package service

import (
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
	"os"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/model"
	"xiet26/face-recognition-local-service/utilities"
)

type AddFace struct {
	FaceID int32 `json:"faceID" bson:"faceID"`
}

func (c *AddFace) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddFaceHandler struct {
	FaceRepository database.FaceMongoRepository
	Recognizer     *face.Recognizer
	RootFolder     string
}

func (h *AddFaceHandler) Handle(c *AddFace, image []byte) error {
	if err := c.Valid(); err != nil {
		return err
	}

	rec, err := face.NewRecognizer(utilities.ModelPath)
	if err != nil {
		return err
	}

	defer rec.Close()

	imageFolder := fmt.Sprintf(utilities.ImagePersonFolderPath, h.RootFolder, c.FaceID)

	faceInfo, err := rec.RecognizeSingle(image)

	if err != nil {
		return err
	}
	if faceInfo == nil {
		return fmt.Errorf("not a single face on the image")
	}

	ok := h.FaceRepository.IsExist(c.FaceID)
	if !ok {
		fmt.Println("create new face")
		e := h.FaceRepository.Create(model.Face{
			ID:     bson.NewObjectId(),
			Vector: faceInfo.Descriptor,
			FaceID: c.FaceID,
		})
		if e != nil {
			return e
		}
	}

	faceData, err := h.FaceRepository.ReadByFaceID(c.FaceID)
	if err != nil {
		return err
	}

	rec.SetSamples([]face.Descriptor{faceData.Vector}, []int32{faceData.FaceID})

	id := rec.ClassifyThreshold(faceInfo.Descriptor, 0.6)

	if id < 0 {
		return fmt.Errorf("can not classify")
	}

	fmt.Println(id)
	if _, err := os.Stat(imageFolder); os.IsNotExist(err) {
		err := os.MkdirAll(imageFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	os.Remove(fmt.Sprintf(utilities.ImageFacePath, imageFolder, c.FaceID))

	f, err := os.OpenFile(fmt.Sprintf(utilities.ImageFacePath, imageFolder, c.FaceID), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(image)
	if err != nil {
		return err
	}

	return nil
}
