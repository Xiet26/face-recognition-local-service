package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/service"
	"xiet26/face-recognition-local-service/utilities"
)

type AttendTempHandler struct {
	FaceRepository database.FaceMongoRepository
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

	if err := handler.Handle(cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Attended"})
}

func (h *AttendTempHandler) GetAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params)  http.Handler{
	batchID := p.ByName("batchID")

	year, _ := GetQuery(r, "year")
	month, _ := GetQuery(r, "month")
	day, _ := GetQuery(r, "day")

	t := fmt.Sprintf("%s-%s-%s", day, month, year)

	filePath := fmt.Sprintf(utilities.ImageBatchFolderPath, batchID, t)

	return http.FileServer(http.Dir(filePath))
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
