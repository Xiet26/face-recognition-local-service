package main

import (
	"flag"
	"fmt"
	"git.cyradar.com/utilities/data/timer"
	"log"
	"net/http"
	"runtime"
	"time"
	"xiet26/face-recognition-local-service/utilities"
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

	//headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	//methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	//origins := handlers.AllowedOrigins([]string{"*"})

	fs := http.FileServer(http.Dir(fmt.Sprintf(`%s/data/images`, container.Config.RootFolder)))
	http.Handle("/", fs)

	container.Logger().Infof("Listen and serve Licence API at %s\n", container.Config.BindingImageService)
	container.Logger().Fatalln(http.ListenAndServe(fmt.Sprintf(`:%s`, container.Config.BindingImageService), nil))

}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&configPrefix, "configPrefix", "studymanagementsystem", "configs prefix")
	flag.StringVar(&configSource, "configSource", "../configs", "configs source")
	flag.StringVar(&modelGoFace, "modelGoFace", "../models", "model for go-face")

}
