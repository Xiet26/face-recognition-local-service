package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"xiet26/face-recognition-local-service/backend/app/api"
)

func NewAPIBeta(container *Container) http.Handler {

	router := api.NewRouter()

	beta := router.Group("/api/beta")

	beta.Use(
		//api.RequireAuth(container.Config.JwtSecret),
		func(handle httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				container.Logger().Infof(`Beta: Method: %s URI:%s`, r.Method, r.RequestURI)
				handle(w, r, p)
			}
		},
	)

	attendRouters(beta)
	faceRouters(beta)
	cameraRouters(beta)

	androidAttendRouters(beta)

	return router
}

func attendRouters(parent *api.Router) {
	handler := api.AttendTempHandler{
		FaceRepository: container.FaceRepository,
		RootFolder:     container.Config.RootFolder,
		ImagePort:      container.Config.BindingImageService,
	}

	router := parent.Group("/attend-temp")
	router.POST("/", handler.AddAttendTemp)
	router.GET("/:batchID/batch", handler.GetAttendTempBatchImages)
	router.GET("/:batchID/face", handler.GetAttendTempFaceImages)
	router.GET("/:batchID/unknown", handler.GetAttendTempFaceImagesUnknown)

	//router.GET("/:batchID", handler.GetAttendTemp)
	//router.DELETE("/:batchID", handler.DeleteAttendTemp)
}

func faceRouters(parent *api.Router) {
	handler := api.FaceHandler{
		FaceRepository: container.FaceRepository,
		RootFolder:     container.Config.RootFolder,
		ImagePort:      container.Config.BindingImageService,
	}

	router := parent.Group("/face")
	router.POST("/", handler.AddFace)
	router.GET("/:faceID", handler.Get)
	router.POST("/faceIDs", handler.GetByFaceIDs)
}

func cameraRouters(parent *api.Router) {
	//handler := &api.CameraHandler{
	//	CameraRepository: container.CameraRepository,
	//}
	//
	//router := parent.Group("/camera")
	//router.GET("", handler.AllCameras)
	//router.PATCH("/:id", handler.UpdateCamera)
	//router.POST("", handler.AddCamera)
	////router.POST("/upload", handler.ImportXLSX)
	//router.DELETE("/:id", handler.DeleteCamera)
}

func androidAttendRouters(parent *api.Router) {
	handler := api.AttendTempHandler{
		FaceRepository: container.FaceRepository,
		RootFolder:     container.Config.RootFolder,
		ImagePort:      container.Config.BindingImageService,
	}

	router := parent.Group("/android/attend-temp")
	router.POST("/", handler.AddAndroidAttendTemp)
	router.GET("/:batchID/batch", handler.GetAndroidAttendTempBatchImages)
}
