package main

import (
	"flag"
	"fmt"
	"git.cyradar.com/utilities/data/timer"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"xiet26/Smart_Attendance_System/backend/utilities"
)

var (
	configPrefix string
	configSource string
	modelGoFace  string

	container *Container
)

func main() {
	flag.Parse()
	defer timer.TimeTrack(time.Now(), fmt.Sprintf("Licence API"))
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			main()
		}
	}()

	var config Config
	err := utilities.LoadEnv(&config, configPrefix, configSource)
	if err != nil {
		log.Fatalln(err)
	}

	container, err = NewContainer(config)
	if err != nil {
		log.Fatalln(err)
	}

	//err = container.GetDataFaceStudent(container.Config.LicenceID)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	os.MkdirAll(container.Config.RootFolder, os.ModePerm)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	container.Logger().Infof("Listen and serve Licence API at %s\n", container.Config.Binding)
	container.Logger().Fatalln(http.ListenAndServe(container.Config.Binding, handlers.CORS(headers, methods, origins)(NewAPIBeta(container))))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&configPrefix, "configPrefix", "studymanagementsystem", "configs prefix")
	flag.StringVar(&configSource, "configSource", "../configs", "configs source")
	flag.StringVar(&modelGoFace, "modelGoFace", "../models", "model for go-face")

}
