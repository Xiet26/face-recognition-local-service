package service

import (
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/asaskevich/govalidator"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
	"time"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/model"
	"xiet26/face-recognition-local-service/utilities"
)

const tolerance = 0.2

type AddAttendTemp struct {
	BatchID string       `json:"batchID"`
	Camera  model.Camera `json:"camera"`
	FaceIDs []int32      `json:"faceIDs"`
}

func (c *AddAttendTemp) Valid() error {
	if c.BatchID == "" {
		return fmt.Errorf("invalid struct")
	}

	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddAttendTempHandler struct {
	FaceRepository database.FaceMongoRepository
	RootFolder     string
}

func (h *AddAttendTempHandler) Handle(data *AddAttendTemp) error {
	if err := data.Valid(); err != nil {
		return err
	}
	fmt.Println(data)
	t := time.Now().Format(utilities.BIRTH_FORMAT_ATTEND)
	folderPath := fmt.Sprintf(utilities.ImageBatchFolderPath, h.RootFolder, data.BatchID, t)

	err := os.MkdirAll(fmt.Sprintf(`%s/all`, folderPath), os.ModePerm)
	if err != nil {
		return err
	}

	imagePaths, err := data.Camera.GetFrames(fmt.Sprintf(`%s/all`, folderPath), time.Now(), 1)
	if err != nil {
		return err
	}

	rec, err := face.NewRecognizer(utilities.ModelPath)
	if err != nil {
		return err
	}

	faceData, err := h.FaceRepository.ReadByMultiFaceID(data.FaceIDs)
	if err != nil {
		return err
	}

	var (
		vectors []face.Descriptor
		ids     []int32
	)

	for _, v := range faceData {
		vectors = append(vectors, v.Vector)
		ids = append(ids, v.FaceID)
	}

	rec.SetSamples(vectors, ids)

	var studentAttendsTmp []model.StudentAttend
	for _, imgPath := range imagePaths {
		facesPath, facesID, e := h.PredictImage(rec, imgPath, data.BatchID, folderPath)
		if e != nil {
			fmt.Println(e)
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

	return nil
}

func (h *AddAttendTempHandler) PredictImage(rec *face.Recognizer, imagePath string, batchID string, path string) ([]string, []int, error) {
	faces, err := rec.RecognizeFile(imagePath)
	if err != nil || faces == nil {
		return nil, nil, fmt.Errorf("can't reconize image")
	}

	folderPath := fmt.Sprintf(`%s/face`, path)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	facePaths := make([]string, 0)
	facesID := make([]int, 0)

	for i, f := range faces {
		id, err := h.predict(rec, f.Descriptor)

		if err != nil {

			drawLineInImage(imagePath, imagePath, f.Rectangle)

			imageFace := fmt.Sprintf(utilities.ImageBatchPath, folderPath, fmt.Sprintf("unknown%d", i), id, time.Now().Unix())

			cropFaceFromImage(imagePath, imageFace, f.Rectangle)

			continue
		}

		drawLineInImage(imagePath, imagePath, f.Rectangle)

		imageFace := fmt.Sprintf(utilities.ImageBatchPath, folderPath, batchID, id, time.Now().Unix())

		cropFaceFromImage(imagePath, imageFace, f.Rectangle)

		facePaths = append(facePaths, imageFace)
		facesID = append(facesID, id)
	}

	// remove image get from camera
	//os.RemoveAll(imagePath)

	return facePaths, facesID, nil
}

func (h *AddAttendTempHandler) predict(rec *face.Recognizer, vector [128]float32) (int, error) {
	id := rec.ClassifyThreshold(vector, tolerance)
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

	//rectangle.Min.X -= 10
	//rectangle.Min.Y -= 10
	//rectangle.Max.X += 10
	//rectangle.Max.Y += 10
	//
	//if rectangle.Min.X < 0 {
	//	rectangle.Min.X = 0
	//}
	//
	//if rectangle.Min.Y < 0 {
	//	rectangle.Min.Y = 0
	//}
	//
	//if rectangle.Max.X > mat.Cols() {
	//	rectangle.Max.X = mat.Cols()
	//}
	//
	//if rectangle.Max.Y > mat.Rows() {
	//	rectangle.Max.Y = mat.Rows()
	//}

	mat = mat.Region(rectangle)
	gocv.IMWrite(dst, mat)
}

func drawLineInImage(src string, dst string, rectangle image.Rectangle) {
	mat := gocv.IMRead(src, gocv.IMReadUnchanged)
	cloneMat := mat.Clone()

	gocv.Rectangle(&cloneMat, rectangle, color.RGBA{G: 255}, 2) // color: green

	gocv.IMWrite(dst, cloneMat)
}
