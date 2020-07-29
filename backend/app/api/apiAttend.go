package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/service"
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

//func (h *AttendTempHandler) GetAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//	batchID := p.ByName("batchID")
//
//	cmd := new(service.GetAttendTemp)
//	cmd.BatchID = batchID
//
//	handler := &service.GetAttendTempHandler{
//		AttendTempRepository: h.AttendTempRepository,
//	}
//
//	result, err := handler.Handle(cmd)
//	if err != nil {
//		ResponseError(w, r, err)
//		return
//	}
//
//	WriteJSON(w, http.StatusOK, ResponseBody{Data: result})
//}
//
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
