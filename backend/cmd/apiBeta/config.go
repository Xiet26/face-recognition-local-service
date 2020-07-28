package main

type Config struct {
	LogLevel   int
	RootFolder string

	// LicenceID
	LicenceID string

	// Binding
	Binding string

	// MongoDB
	MongoServer     string
	MongoDatabase   string
	MongoCollection string
	MongoSource     string
	MongoUser       string
	MongoPassword   string
}
