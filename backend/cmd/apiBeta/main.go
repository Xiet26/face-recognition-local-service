package main

import (
	"flag"
	"fmt"
	"git.cyradar.com/utilities/data/timer"
	"github.com/Kagami/go-face"
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

	container.Recognizer, err = face.NewRecognizer(modelGoFace)
	if err != nil {
		log.Fatalln(err)
	}

	container.GetDataFaceStudent(container.Config.LicenceID)
	//err = container.GetDataFaceStudent(container.Config.LicenceID)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	os.MkdirAll(container.Config.RootFolder, os.ModePerm)

	container.Logger().Infof("Listen and serve Licence API at %s\n", container.Config.Binding)
	container.Logger().Fatalln(http.ListenAndServe(container.Config.Binding, NewAPIBeta(container)))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&configPrefix, "configPrefix", "studymanagementsystem", "configs prefix")
	flag.StringVar(&configSource, "configSource", "../configs", "configs source")
	flag.StringVar(&modelGoFace, "modelGoFace", "../models-go-face", "model for go-face")

}
