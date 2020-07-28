package api

import (
	"github.com/Kagami/go-face"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"xiet26/goface/face-management/database"
	"xiet26/goface/face-management/service"
)

type AttendTempHandler struct {
	AttendTempRepository database.AttendTempMongoRepository
	Recognizer           *face.Recognizer
	RootFolder           string
}

func (h *AttendTempHandler) AddAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cmd := new(service.AddAttendTemp)
	if err := BindJSON(r, cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	handler := &service.AddAttendTempHandler{
		AttendTempRepository: h.AttendTempRepository,
		Recognizer:           h.Recognizer,
		RootFolder:           h.RootFolder,
	}

	if err := handler.Handle(cmd); err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Attended"})
}

func (h *AttendTempHandler) GetAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	batchID := p.ByName("batchID")

	cmd := new(service.GetAttendTemp)
	cmd.BatchID = batchID

	handler := &service.GetAttendTempHandler{
		AttendTempRepository: h.AttendTempRepository,
	}

	result, err := handler.Handle(cmd)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Data: result})
}

func (h *AttendTempHandler) DeleteAttendTemp(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	batchID := p.ByName("batchID")

	cmd := new(service.DeleteAttendTemp)
	cmd.BatchID = batchID

	handler := &service.DeleteAttendTempHandler{
		AttendTempRepository: h.AttendTempRepository,
	}

	err := handler.Handle(cmd)
	if err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Successful"})
}
