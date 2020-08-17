package api

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/julienschmidt/httprouter"
	"github.com/rwcarlsen/goexif/exif"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"strconv"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/service"
	"xiet26/face-recognition-local-service/utilities"
)

type FaceHandler struct {
	FaceRepository database.FaceMongoRepository
	ImagePort      string
	RootFolder     string
}

func (h *FaceHandler) AddFace(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(service.AddFace)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	f, _, err := r.FormFile(utilities.FileFieldFaceImage)
	if err != nil {
		ResponseError(w, r, err)
		return
	}
	defer f.Close()

	image, err := ioutil.ReadAll(f)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	faceId := r.FormValue(utilities.FaceIDField)

	id, err := strconv.Atoi(faceId)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	cmd.FaceID = int32(id)
	handler := &service.AddFaceHandler{
		FaceRepository: h.FaceRepository,
		RootFolder:     h.RootFolder,
	}

	if err := handler.Handle(cmd, image); err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Attended"})
}

func (h *FaceHandler) AndroidAddFace(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	f, _, err := r.FormFile(utilities.FileFieldFaceImage)
	if err != nil {
		ResponseError(w, r, err)
		return
	}
	defer f.Close()

	tmpImageByte, err := ioutil.ReadAll(f)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	var tmpReader01 = bytes.NewReader(tmpImageByte)
	var tmpReader02 = bytes.NewReader(tmpImageByte)

	tmpImage, err := jpeg.Decode(tmpReader01)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	x, err := exif.Decode(tmpReader02)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	var rotation float64 = 0
	orientationRaw, err := x.Get("Orientation")
	if err == nil {
		orientation := orientationRaw.String()
		if orientation == "3" {
			rotation = 180
		} else if orientation == "6" {
			rotation = 270
		} else if orientation == "8" {
			rotation = 90
		}
	}

	img := imaging.Rotate(tmpImage, rotation, color.Gray{})

	var buf = new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	image := buf.Bytes()

	faceID := p.ByName("faceID")
	id, err := strconv.Atoi(faceID)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	cmd := new(service.AddFace)
	cmd.FaceID = int32(id)
	handler := &service.AddFaceHandler{
		FaceRepository: h.FaceRepository,
		RootFolder:     h.RootFolder,
	}

	if err := handler.Handle(cmd, image); err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Message: "added face"})
}

func (h *FaceHandler) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	faceID := p.ByName("faceID")

	hostPath := fmt.Sprintf(`:%s/person/%s`, h.ImagePort, faceID)
	folderPath := fmt.Sprintf(utilities.ImagePersonFolderPath, h.RootFolder, faceID)

	var filePaths []string

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	for _, f := range files {
		filePaths = append(filePaths, fmt.Sprintf(`%s/%s`, hostPath, f.Name()))
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "get face images",
		Code:    0,
		Data:    filePaths,
	})

	return
}

func (h *FaceHandler) GetByFaceIDs(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	type Body struct {
		FaceIDs []int `json:"faceIDs"`
	}
	var req = new(Body)
	if err := BindJSON(r, req); err != nil {
		ResponseError(w, r, err)
		return
	}

	type ResponseData struct {
		FaceID int32  `json:"faceID"`
		URL    string `json:"url"`
	}
	var result = make([]ResponseData, 0)

	for _, faceID := range req.FaceIDs {
		hostPath := fmt.Sprintf(`:%s/person/%d`, h.ImagePort, faceID)
		folderPath := fmt.Sprintf(utilities.ImagePersonFolderPath, h.RootFolder, faceID)

		files, err := ioutil.ReadDir(folderPath)
		if err != nil || len(files) < 1 {
			continue
		}

		result = append(result, ResponseData{
			FaceID: int32(faceID),
			URL:    fmt.Sprintf(`%s/%s`, hostPath, files[0].Name()),
		})
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "get face images",
		Code:    0,
		Data:    result,
	})

	return
}
