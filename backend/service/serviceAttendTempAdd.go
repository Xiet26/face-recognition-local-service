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
	Group   int          `json:"group"`
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

func (h *AddAttendTempHandler) Handle(data *AddAttendTemp) (model.ResultAttend, error, string) {
	errMsg := "Error in time: "
	var result model.ResultAttend
	if err := data.Valid(); err != nil {
		return result, err, errMsg
	}

	timeNow := time.Now()
	fmt.Println(data)

	t := timeNow.Format(utilities.BIRTH_FORMAT_ATTEND)
	folderPath := fmt.Sprintf(utilities.ImageBatchFolderPath, h.RootFolder, data.BatchID, data.Group, t)
	err := os.MkdirAll(fmt.Sprintf(`%s/all`, folderPath), os.ModePerm)
	if err != nil {
		return result, err, errMsg
	}

	imagePaths, err := data.Camera.GetFrames(folderPath, time.Now(), 1)
	if err != nil {
		return result, err, errMsg
	}

	rec, err := face.NewRecognizer(utilities.ModelPath)
	if err != nil {
		return result, err, errMsg
	}

	faceData, err := h.FaceRepository.ReadByMultiFaceID(data.FaceIDs)
	if err != nil {
		return result, err, errMsg
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

	studentAttends := make(map[int]bool)

	for _, imgPath := range imagePaths {
		_, facesIDs, e := h.PredictImage(rec, imgPath, data.BatchID, folderPath, timeNow.Unix())
		if e != nil {
			errMsg = fmt.Sprintf("%v", timeNow)
			continue
		}

		for _, v := range facesIDs {
			if studentAttends[v] {
				continue
			}
			studentAttends[v] = true
		}
	}

	var absent []int32
	for _, v := range data.FaceIDs {
		if !studentAttends[int(v)] {
			absent = append(absent, v)
		}
	}

	result.Time = timeNow
	result.BatchID = data.BatchID
	result.Group = data.Group
	result.AbsentFaceIDs = absent

	return result, err, errMsg
}

func (h *AddAttendTempHandler) AndroidHandle(data *AddAttendTemp) (model.ResultAndroidAttend, error, string) {
	errMsg := "Error in time: "
	var result model.ResultAndroidAttend
	if err := data.Valid(); err != nil {
		return result, err, errMsg
	}

	timeNow := time.Now()
	fmt.Println(timeNow.Unix())

	t := timeNow.Format(utilities.BIRTH_FORMAT_ATTEND)
	folderPath := fmt.Sprintf(utilities.ImageBatchFolderPath, h.RootFolder, data.BatchID, data.Group, t)

	err := os.MkdirAll(fmt.Sprintf(`%s/all`, folderPath), os.ModePerm)
	if err != nil {
		return result, err, errMsg
	}

	imagePaths, err := data.Camera.GetFrames(folderPath, time.Now(), 1)
	if err != nil {
		return result, err, errMsg
	}

	rec, err := face.NewRecognizer(utilities.ModelPath)
	if err != nil {
		return result, err, errMsg
	}

	faceData, err := h.FaceRepository.ReadByMultiFaceID(data.FaceIDs)
	if err != nil {
		return result, err, errMsg
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

	studentAttends := make(map[int]bool)

	for _, imgPath := range imagePaths {
		_, facesIDs, e := h.PredictImage(rec, imgPath, data.BatchID, folderPath, timeNow.Unix())
		if e != nil {
			errMsg = fmt.Sprintf("%v", timeNow)
			continue
		}

		for _, v := range facesIDs {
			if studentAttends[v] {
				continue
			}
			studentAttends[v] = true
		}
	}

	var absent []int32
	for _, v := range data.FaceIDs {
		if !studentAttends[int(v)] {
			absent = append(absent, v)
		}
	}

	result.Time = timeNow.Format("02-01-2006 15:04:05")
	result.BatchID = data.BatchID
	result.Group = data.Group
	result.AbsentFaceIDs = absent

	return result, err, errMsg
}

func (h *AddAttendTempHandler) PredictImage(rec *face.Recognizer, imagePath string, batchID string, path string, t int64) ([]string, []int, error) {
	defer os.RemoveAll(imagePath)
	folderAllPath := fmt.Sprintf(`%s/all`, path)

	faces, err := rec.RecognizeFile(imagePath)
	if err != nil || faces == nil {
		imageAll := fmt.Sprintf(utilities.ImageBatchAllPath, folderAllPath, t)
		drawLineInImage(imagePath, imageAll, image.Rectangle{})
		return nil, nil, fmt.Errorf("can't reconize image")
	}

	folderFacePath := fmt.Sprintf(`%s/face`, path)
	err = os.MkdirAll(folderFacePath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	err = os.MkdirAll(folderAllPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	var facePaths []string
	var faceIDs []int

	for i, f := range faces {
		id, err := h.predict(rec, f.Descriptor)

		imageAll := fmt.Sprintf(utilities.ImageBatchAllPath, folderAllPath, t)

		drawLineInImage(imagePath, imageAll, f.Rectangle)

		if err != nil {
			imageFace := fmt.Sprintf(utilities.ImageBatchPath, folderFacePath, fmt.Sprintf("unknown%d", i), id, t)

			cropFaceFromImage(imagePath, imageFace, f.Rectangle)

			continue
		}

		imageFace := fmt.Sprintf(utilities.ImageBatchPath, folderFacePath, batchID, id, t)

		cropFaceFromImage(imagePath, imageFace, f.Rectangle)

		facePaths = append(facePaths, imageFace)
		faceIDs = append(faceIDs, id)
	}

	return facePaths, faceIDs, nil
}

func (h *AddAttendTempHandler) predict(rec *face.Recognizer, vector [128]float32) (int, error) {
	id := rec.ClassifyThreshold(vector, tolerance)
	if id < 0 {
		return -1, fmt.Errorf("cant classify")
	}

	return id, nil
}


func cropFaceFromImage(src string, dst string, rectangle image.Rectangle) {
	mat := gocv.IMRead(src, gocv.IMReadUnchanged)

	rateX := float64(rectangle.Dx())/float64(mat.Cols())
	rateY := float64(rectangle.Dy())/float64(mat.Rows())

	rectangle.Min.X -= int(200*rateX)
	rectangle.Min.Y -= int(200*rateY)
	rectangle.Max.X += int(150*rateX)
	rectangle.Max.Y += int(150*rateX)

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

func drawLineInImage(src string, dst string, rectangle image.Rectangle) {
	mat := gocv.IMRead(src, gocv.IMReadUnchanged)
	cloneMat := mat.Clone()

	gocv.Rectangle(&cloneMat, rectangle, color.RGBA{R: 255}, 1) // color: green

	gocv.IMWrite(dst, cloneMat)

}
