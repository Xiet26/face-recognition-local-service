package main

import (
	"net/http"
	"xiet26/goface/face-management/app/api"
)

func NewAPIBeta(container *Container) http.Handler {

	router := api.NewRouter()

	beta := router.Group("/api/beta")

	//beta.Use(
	//	api.RequireAuth(container.Config.JwtSecret),
	//	func(handle httprouter.Handle) httprouter.Handle {
	//		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//			container.Logger().Infof(`Beta: Method: %s URI:%s`, r.Method, r.RequestURI)
	//			handle(w, r, p)
	//		}
	//	},
	//)

	attendRouters(beta)
	cameraRouters(beta)

	return router
}

func attendRouters(parent *api.Router) {
	handler := api.AttendTempHandler{
		AttendTempRepository: container.AttendTempRepository,
		Recognizer:           container.Recognizer,
		RootFolder:           container.Config.RootFolder,
	}

	router := parent.Group("/attend-temp")
	router.GET("/xxx", handler.Test)
	//router.GET("/", handler.AddAttendTemp)
	//router.GET("/:batchID", handler.GetAttendTemp)
	//router.DELETE("/:batchID", handler.DeleteAttendTemp)
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
