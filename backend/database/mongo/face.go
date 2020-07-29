package mongoDatabase

import (
	"git.cyradar.com/utilities/data/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"xiet26/face-recognition-local-service/backend/model"
)

var FaceMongoCollection = "faces"

type FaceMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewFaceMongoRepository(provider *mongo.MongoProvider) *FaceMongoRepository {
	repo := &FaceMongoRepository{provider, FaceMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"faceID",
		},
		Unique: true,
	})

	return repo
}

func (repo *FaceMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *FaceMongoRepository) Create(face model.Face) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(face))
}

func (repo *FaceMongoRepository) ReadByFaceID(faceID int32) (model.Face, error) {
	collection, close := repo.collection()
	defer close()

	var result model.Face
	err := collection.Find(bson.M{
		"faceID": faceID,
	}).One(&result)
	return result, repo.provider.NewError(err)
}

func (repo *FaceMongoRepository) ReadByMultiFaceID(faceIDs []int32) ([]model.Face, error) {
	collection, close := repo.collection()
	defer close()

	var result []model.Face
	err := collection.Find(bson.M{
		"faceID": bson.M{"$in": faceIDs},
	}).One(&result)
	return result, repo.provider.NewError(err)
}

func (repo *FaceMongoRepository) DeleteByFaceID(faceID int32) error {
	collection, close := repo.collection()
	defer close()

	_, err := collection.RemoveAll(bson.M{
		"faceID": faceID,
	})

	return repo.provider.NewError(err)
}

func (repo *FaceMongoRepository) IsExist(faceID int32) bool {
	collection, close := repo.collection()
	defer close()

	n, err := collection.Find(bson.M{
		"faceID": faceID,
	}).Count()

	if err != nil || n == 0 {
		return false
	}
	return true
}
