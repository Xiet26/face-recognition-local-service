package service

import (
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/asaskevich/govalidator"
	"gocv.io/x/gocv"
	"image"
	"os"
	"time"
	"xiet26/goface/face-management/database"
	"xiet26/goface/face-management/model"
)

const tolerance = 0.2

type AddAttendTemp struct {
	BatchID    string    `json:"batchID"`
	Time       time.Time `json:"time"`
	CameraHost string    `json:"cameraHost"`
	CameraPort string    `json:"cameraPort"`
}

func (c *AddAttendTemp) Valid() error {
	if c.BatchID == "" {
		return fmt.Errorf("invalid struct")
	}

	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddAttendTempHandler struct {
	AttendTempRepository database.AttendTempMongoRepository
	Recognizer           *face.Recognizer
	RootFolder           string
}

func (h *AddAttendTempHandler) Handle(data *AddAttendTemp) error {
	if err := data.Valid(); err != nil {
		return err
	}

	cam := model.Camera{
		Host: data.CameraHost,
		Port: data.CameraPort,
	}

	imagePaths, err := cam.GetFrames(h.RootFolder, data.Time, 5)
	if err != nil {
		return err
	}

	var studentAttendsTmp []model.StudentAttend
	for _, imgPath := range imagePaths {
		facesPath, facesID, e := h.PredictImage(imgPath, data.BatchID)
		if e != nil {
			continue
		}

		for i := 0; i < len(facesPath); i++ {
			studentAttendsTmp = append(studentAttendsTmp, model.StudentAttend{
				ImageFace: facesPath[i],
				FaceID:    facesID[i],
			})
		}
	}

	var studentAttends []model.StudentAttend
	for _, v := range studentAttendsTmp {
		if isExistedStudent(v.FaceID, studentAttends) {
			continue
		}

		studentAttends = append(studentAttends, v)
	}

	return h.AttendTempRepository.Create(model.AttendTemp{
		BatchID:        data.BatchID,
		Time:           data.Time,
		StudentAttends: studentAttends,
	})
}

func (h *AddAttendTempHandler) PredictImage(imagePath string, batchID string) ([]string, []int, error) {
	faces, err := h.Recognizer.RecognizeFile(imagePath)
	if err != nil || faces == nil {
		return nil, nil, fmt.Errorf("can't reconize image")
	}

	os.MkdirAll(fmt.Sprintf("%s/%s", h.RootFolder, batchID), os.ModePerm)

	facePaths := make([]string, 0)
	facesID := make([]int, 0)

	for _, f := range faces {
		id, err := h.predict(f.Descriptor)
		if err != nil {
			continue
		}

		imageFace := fmt.Sprintf("%s/%s/%v.png", h.RootFolder, batchID, time.Now().Unix())
		cropFaceFromImage(imagePath, imageFace, f.Rectangle)

		facePaths = append(facePaths, imageFace)
		facesID = append(facesID, id)
	}

	// remove image get from camera
	os.RemoveAll(imagePath)

	return facePaths, facesID, nil
}

func (h *AddAttendTempHandler) predict(vector [128]float32) (int, error) {
	id := h.Recognizer.ClassifyThreshold(vector, tolerance)
	if id < 0 {
		return -1, fmt.Errorf("cant classify")
	}

	return id, nil
}

func isExistedStudent(id int, studentAttends []model.StudentAttend) bool {
	for _, v := range studentAttends {
		if v.FaceID == id {
			return true
		}
	}
	return false
}

func cropFaceFromImage(src string, dst string, rectangle image.Rectangle) {
	mat := gocv.IMRead(src, gocv.IMReadUnchanged)

	rectangle.Min.X -= 150
	rectangle.Min.Y -= 150
	rectangle.Max.X += 100
	rectangle.Max.Y += 100

	if rectangle.Min.X < 0 {
		rectangle.Min.X = 0
	}

	if rectangle.Min.Y < 0 {
		rectangle.Min.Y = 0
	}

	if rectangle.Max.X > mat.Cols() {
		rectangle.Max.X = mat.Cols()
	}

	if rectangle.Max.Y > mat.Rows() {
		rectangle.Max.Y = mat.Rows()
	}

	mat = mat.Region(rectangle)
	gocv.IMWrite(dst, mat)
}
