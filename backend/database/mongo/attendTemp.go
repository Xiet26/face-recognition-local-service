package mongoDatabase

import (
	"git.cyradar.com/utilities/data/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"xiet26/goface/face-management/model"
)

var AttendTempMongoCollection = "attendTemps"

type AttendTempMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewAttendTempMongoRepository(provider *mongo.MongoProvider) *AttendTempMongoRepository {
	repo := &AttendTempMongoRepository{provider, AttendTempMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"batchID",
		},
	})

	return repo
}

func (repo *AttendTempMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *AttendTempMongoRepository) Create(temp model.AttendTemp) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(temp))
}

func (repo *AttendTempMongoRepository) ReadByBatchID(batchID string) ([]model.AttendTemp, error) {
	collection, close := repo.collection()
	defer close()

	var result []model.AttendTemp
	err := collection.Find(bson.M{
		"batchID":  batchID,
	}).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *AttendTempMongoRepository) DeleteByBatchID(batchID string) error {
	collection, close := repo.collection()
	defer close()

	_, err := collection.RemoveAll(bson.M{
		"batchID":  batchID,
	})

	return repo.provider.NewError(err)
}
