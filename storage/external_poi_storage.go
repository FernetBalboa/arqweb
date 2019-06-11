package storage

import (
	"context"
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

//go:generate mockgen -destination=../mock/mock_external_poi_storage.go -package=mock -source=external_poi_storage.go -imports geojson=github.com/paulmach/go.geojson
type ExternalPOIStorage interface {
	SavePOI(POI *domain.ExternalPOI) (*domain.ExternalPOI, error)
	GetPOIs() ([]domain.ExternalPOI, error)
	GetPOI(poiID primitive.ObjectID) (*domain.ExternalPOI, error)
	RemovePOI(poiID primitive.ObjectID) (*mongo.DeleteResult, error)
}

const (
	ExternalPOICollection = "external_poi"
)

func init() {
	//resetPoiCollection() //Comment if data should be kept between program runs
}

type ExternalPOIStorageImpl struct {
	poiCollection ICollection
}

func CreateExternalPOIStorage(POIcollection ICollection) (ExternalPOIStorage, error) {
	storage := &ExternalPOIStorageImpl{
		poiCollection: POIcollection,
	}

	return storage, nil
}

func NewExternalPOIStorage() (ExternalPOIStorage, error) {
	client, err := getMongoDBClient()
	poiCollection := client.Database(Database).Collection(ExternalPOICollection)
	if err != nil {
		return nil, err
	}

	return CreateExternalPOIStorage(poiCollection)
}

func (ps *ExternalPOIStorageImpl) SavePOI(POI *domain.ExternalPOI) (*domain.ExternalPOI, error) {

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := ps.poiCollection.InsertOne(ctx, POI)

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not insert new external POI into MongoDB. POI: %+v", POI)
	}

	POI.Id = res.InsertedID.(primitive.ObjectID)

	return POI, nil
}

func (ps *ExternalPOIStorageImpl) GetPOIs() ([]domain.ExternalPOI, error) {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	resCursor, err := ps.poiCollection.Find(ctx, bson.M{})

	if err != nil {
		return []domain.ExternalPOI{}, apierror.Wrapf(err, "Could not retrieve list of external POIs")
	}

	defer resCursor.Close(ctx)

	results := []domain.ExternalPOI{}
	for resCursor.Next(ctx) {
		var poi domain.ExternalPOI
		err = resCursor.Decode(&poi)
		if err != nil {
			log.Errorf("Error while decoding external POI. Cause: %v", err)
		} else {
			results = append(results, poi)
		}
	}

	return results, nil
}

func (ps *ExternalPOIStorageImpl) GetPOI(poiID primitive.ObjectID) (*domain.ExternalPOI, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	log.Infof("Getting external POI '%s'", poiID)

	// set filters
	filter := bson.M{"_id": poiID}

	var poi domain.ExternalPOI
	err := ps.poiCollection.FindOne(ctx, filter).Decode(&poi)

	if err != nil {
		return nil, apierror.Wrapf(err, "Couldn't get external POI with id: '%s'", poiID)
	}

	return &poi, nil
}

func (ps *ExternalPOIStorageImpl) RemovePOI(poiID primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := ps.poiCollection.DeleteOne(ctx, bson.M{"_id": poiID})

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not remove external POI from MongoDB. POI: %s", poiID)
	}

	return res, nil
}
