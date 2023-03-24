package main

import (
	"git.cyradar.com/utilities/data/providers/logger"
	mongoDatabase "xiet26/face-recognition-local-service/backend/database/mongo"
)

type Providers struct {
	*logger.LoggerProvider
}

type Container struct {
	*Providers

	AttendTempRepository *mongoDatabase.AttendTempMongoRepository
	FaceRepository       *mongoDatabase.FaceMongoRepository
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


	return nil
}

func (container *Container) LoadDependencies(config Config) {
	container.Config = config
}

func (container *Container) LoadProviders(config Config) {

	container.Providers = &Providers{
		LoggerProvider: logger.NewLoggerProvider(config.LogLevel),
	}
}

