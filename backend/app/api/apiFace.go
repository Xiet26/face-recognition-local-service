package api

import (
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strconv"
	"xiet26/face-recognition-local-service/backend/database"
	"xiet26/face-recognition-local-service/backend/service"
	"xiet26/face-recognition-local-service/utilities"
)

type FaceHandler struct {
	FaceRepository database.FaceMongoRepository
	Recognizer     *face.Recognizer
	RootFolder     string
}

func (h *FaceHandler) AddFace(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("ASDASDASDASD")
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
		Recognizer:     h.Recognizer,
		RootFolder:     h.RootFolder,
	}

	if err := handler.Handle(cmd, image); err != nil {
		ResponseError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Attended"})
}

func (h *FaceHandler) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("ASDASDASDASD")
	WriteJSON(w, http.StatusOK, ResponseBody{Message: "Attended"})
}
