package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/service"
	"xiet26/face-recognition-local-service/utilities"
)

type AttendTempHandler struct {
	FaceRepository database.FaceMongoRepository
	ImagePort      string
	RootFolder     string
}

func (h *AttendTempHandler) AddAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(service.AddAttendTemp)
	if err := BindJSON(r, cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	handler := &service.AddAttendTempHandler{
		FaceRepository: h.FaceRepository,
		RootFolder:     h.RootFolder,
	}

	result, err, message := handler.Handle(cmd)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: message,
		Code:    0,
		Data:    result,
	})
}

func (h *AttendTempHandler) AddAndroidAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(service.AddAttendTemp)
	if err := BindJSON(r, cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	handler := &service.AddAttendTempHandler{
		FaceRepository: h.FaceRepository,
		RootFolder:     h.RootFolder,
	}

	result, err, message := handler.AndroidHandle(cmd)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: message,
		Code:    0,
		Data:    result,
	})
}

func (h *AttendTempHandler) GetAttendTempBatchImages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	batchID := p.ByName("batchID")

	year, _ := GetQuery(r, "year")
	month, _ := GetQuery(r, "month")
	day, _ := GetQuery(r, "day")
	group, _ := GetQuery(r, "group")

	g, err := strconv.Atoi(group)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	t := fmt.Sprintf("%s-%s-%s", day, month, year)

	hostPath := fmt.Sprintf(`:%s/batch/%s/%v/%s/all`, h.ImagePort, batchID, group, t)
	folderPath := fmt.Sprintf(utilities.ImageBatchFolderPathAll, h.RootFolder, batchID, g, t)

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
		Message: "get batch images",
		Code:    0,
		Data:    filePaths,
	})
}

func (h *AttendTempHandler) GetAttendTempFaceImages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	batchID := p.ByName("batchID")

	year, _ := GetQuery(r, "year")
	month, _ := GetQuery(r, "month")
	day, _ := GetQuery(r, "day")
	group, _ := GetQuery(r, "group")
	faceID, _ := GetQuery(r, "faceID")

	t := fmt.Sprintf("%s-%s-%s", day, month, year)

	g, err := strconv.Atoi(group)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	folderPath := fmt.Sprintf(utilities.ImageBatchFolderPathFace, h.RootFolder, batchID, g, t)
	hostPath := fmt.Sprintf(`:%s/batch/%s/%v/%s/face`, h.ImagePort, batchID, group, t)

	var filePaths []string

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	for _, f := range files {
		if strings.Contains(f.Name(), fmt.Sprintf(`face_%s_`, faceID)) {
			filePaths = append(filePaths, fmt.Sprintf(`%s/%s`, hostPath, f.Name()))
		}
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "get face images",
		Code:    0,
		Data:    filePaths,
	})
	return
}

func (h *AttendTempHandler) GetAttendTempFaceImagesUnknown(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	batchID := p.ByName("batchID")

	year, _ := GetQuery(r, "year")
	month, _ := GetQuery(r, "month")
	day, _ := GetQuery(r, "day")
	group, _ := GetQuery(r, "group")

	t := fmt.Sprintf("%s-%s-%s", day, month, year)

	g, err := strconv.Atoi(group)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	folderPath := fmt.Sprintf(utilities.ImageBatchFolderPathFace, h.RootFolder, batchID, g, t)
	hostPath := fmt.Sprintf(`:%s/batch/%s/%v/%s/face`, h.ImagePort, batchID, group, t)

	var filePaths []string

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	for _, f := range files {
		if strings.Contains(f.Name(), fmt.Sprintf(`face_-1_`)) {
			filePaths = append(filePaths, fmt.Sprintf(`%s/%s`, hostPath, f.Name()))
		}
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "get face images",
		Code:    0,
		Data:    filePaths,
	})
	return
}

//func (h *AttendTempHandler) DeleteAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	batchID := p.ByName("batchID")
//
//	cmd := new(service.DeleteAttendTemp)
//	cmd.BatchID = batchID
//
//	handler := &service.DeleteAttendTempHandler{
//		AttendTempRepository: h.AttendTempRepository,
//	}
//
//	err := handler.Handle(cmd)
//	if err != nil {
//		ResponseError(w, r, err)
//		return
//	}
//
//	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Successful"})
//}