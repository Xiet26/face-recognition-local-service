package main

import (
	"encoding/json"
	"fmt"
	"git.cyradar.com/utilities/data/providers/logger"
	"git.cyradar.com/utilities/data/providers/mongo"
	"github.com/Kagami/go-face"
	"io/ioutil"
	"net/http"
	mongoDatabase "xiet26/goface/face-management/database/mongo"
)

type Providers struct {
	*logger.LoggerProvider
	*mongo.MongoProvider
}

type Container struct {
	*Providers

	Recognizer           *face.Recognizer
	AttendTempRepository *mongoDatabase.AttendTempMongoRepository
	Config               Config
}

func NewContainer(config Config) (*Container, error) {
	container := new(Container)
	err := container.InitContainer(config)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (container *Container) InitContainer(config Config) error {
	// load dependencies
	container.LoadDependencies(config)

	// Load providers into container
	container.LoadProviders(config)

	// Load repositories
	container.LoadRepositoryImplementations()

	return nil
}

func (container *Container) LoadDependencies(config Config) {
	container.Config = config
}

func (container *Container) LoadProviders(config Config) {

	container.Providers = &Providers{
		LoggerProvider: logger.NewLoggerProvider(config.LogLevel),
		MongoProvider: mongo.NewMongoProvider(
			config.MongoServer, config.MongoUser, config.MongoPassword,
			config.MongoDatabase, config.MongoCollection, config.MongoSource,
		),
	}
}

func (container *Container) LoadRepositoryImplementations() {
	container.AttendTempRepository = mongoDatabase.NewAttendTempMongoRepository(container.MongoProvider)
}

func (container *Container) GetDataFaceStudent(licenceID string) error {
	const apiFaceStudent = "http://api/get/face-data"
	resp, err := http.Get(fmt.Sprintf("%s?licenceID=%s", apiFaceStudent, licenceID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type FaceData struct {
		Vector [128]float32 `json:"vector"`
		FaceID int32        `json:"faceID"`
	}

	result := make([]FaceData, 0)
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	var vectors []face.Descriptor
	var faceID []int32

	for _, v := range result {
		vectors = append(vectors, v.Vector)
		faceID = append(faceID, v.FaceID)
	}

	container.Recognizer.SetSamples(vectors, faceID)

	return nil
}
