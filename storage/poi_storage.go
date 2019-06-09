package storage

import (
	"context"
	"github.com/fernetbalboa/arqweb/apierror"
	"github.com/fernetbalboa/arqweb/domain"
	"github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//go:generate mockgen -destination=../mock/mock_poi_storage.go -package=mock -source=poi_storage.go -imports geojson=github.com/paulmach/go.geojson
type POIStorage interface {
	SavePOI(POI *domain.PointOfInterest) (*domain.PointOfInterest, error)
	EditPOI(poiID primitive.ObjectID, newVersionPOI *domain.PointOfInterest) error
	SaveFeature(feature *geojson.Feature) (*domain.PointOfInterest, error)
	SearchPOI(filters *domain.POIFilter) ([]*domain.PointOfInterest, error)
}

const (
	POICollection = "poi"
)

func init() {
	resetPoiCollection() //Comment if data should be kept between program runs
}

type POIStorageImpl struct {
	poiCollection ICollection
}

func CreatePOIStorage(POIcollection ICollection) (POIStorage, error) {
	storage := &POIStorageImpl{
		poiCollection: POIcollection,
	}

	return storage, nil
}

func NewPOIStorage() (POIStorage, error) {
	client, err := getMongoDBClient()
	poiCollection := client.Database(Database).Collection(POICollection)
	if err != nil {
		return nil, err
	}

	return CreatePOIStorage(poiCollection)
}

func (ps *POIStorageImpl) SavePOI(POI *domain.PointOfInterest) (*domain.PointOfInterest, error) {

	POI.Id = primitive.NewObjectID()
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := ps.poiCollection.InsertOne(ctx, POI)

	if err != nil {
		return nil, apierror.Wrapf(err, "Could not insert new POI into MongoDB. POI: %+v", POI)
	}

	POI.Id = res.InsertedID.(primitive.ObjectID)

	return POI, nil
}

func (ps *POIStorageImpl) EditPOI(poiID primitive.ObjectID, newVersionPOI *domain.PointOfInterest) error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	log.Infof("Updating POI '%s'", newVersionPOI.Id)

	// set filters and updates
	filter := bson.M{"_id": poiID}
	update := bson.M{"$set": newVersionPOI}

	// update document
	_, err := ps.poiCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return apierror.Wrapf(err, "Couldn't update POI with id: '%s'", newVersionPOI.Id)
	}

	log.Infof("POI '%s' successfully updated", newVersionPOI.Id)
	return nil
}

func (ps *POIStorageImpl) SaveFeature(feature *geojson.Feature) (*domain.PointOfInterest, error) {
	POI, err := ps.featureToPOI(feature)
	if err != nil {
		return nil, apierror.Wrapf(err,
			"Could not convert GeoJson feature to POI. Feature: %+v", feature)
	}

	return ps.SavePOI(POI)
}

func (ps *POIStorageImpl) featureToPOI(feature *geojson.Feature) (*domain.PointOfInterest, error) {
	title, err := feature.PropertyString("title")
	if err != nil {
		return nil, apierror.BadRequest.Wrapf(err,
			"GeoJson feature 'title' property is missing or has incorrect type (should be string)")
	}

	category, err := feature.PropertyString("category")
	if err != nil {
		return nil, apierror.BadRequest.Wrapf(err,
			"GeoJson feature 'category' property is missing or has incorrect type (should be string)")
	}

	//Not required
	description := feature.PropertyMustString("description", "")

	hidden := feature.PropertyMustBool("hidden", false)

	picture := domain.Picture{}
	if pictureInfo, ok := feature.Properties["picture"].(map[string]interface{}); ok {
		if data, ok2 := pictureInfo["data"].(string); ok2 {
			picture.Data = data
		}
		if contentType, ok2 := pictureInfo["contentType"].(string); ok2 {
			picture.ContentType = contentType
		}
		if name, ok2 := pictureInfo["name"].(string); ok2 {
			picture.Name = name
		}
	}

	if feature.Geometry.IsPoint() {
		lat := feature.Geometry.Point[0]
		long := feature.Geometry.Point[1]
		//long
		POI := &domain.PointOfInterest{
			Title:       title,
			Category:    category,
			Description: description,
			Picture:     picture,
			Lat:         lat,
			Long:        long,
			Hidden:      hidden,
		}
		return POI, nil
	}

	return nil, apierror.BadRequest.New("Only GeoJson Point feature is currently supported")

}

func resetPoiCollection() {
	client, err := getMongoDBClient()

	if err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Database(Database).Collection(POICollection).Drop(ctx)

	if err != nil {
		log.Error("Could not reset POI MongoDB collection")
	}
}

func (ps *POIStorageImpl) SearchPOI(filters *domain.POIFilter) ([]*domain.PointOfInterest, error) {
	if filters == nil {
		filters = &domain.POIFilter{}
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	findOptions := options.Find()
	findOptions.SetLimit(filters.Limit)

	queryFilters := buildPOIQueryFilters(filters)
	resCursor, err := ps.poiCollection.Find(ctx, queryFilters, findOptions)

	if err != nil {
		return []*domain.PointOfInterest{}, apierror.Wrapf(err, "Could not find POIs using filters %+v", filters)
	}

	defer resCursor.Close(ctx)

	//noinspection GoPreferNilSlice
	results := []*domain.PointOfInterest{}
	for resCursor.Next(ctx) {
		var POI domain.PointOfInterest
		err = resCursor.Decode(&POI)
		if err != nil {
			log.Errorf("Error while decoding POI. Cause: %v", err)
		} else {
			results = append(results, &POI)
		}
	}

	return results, nil

}

func buildPOIQueryFilters(filters *domain.POIFilter) bson.M {
	filtersMap := bson.M{}
	if filters.Category != "" {
		filtersMap["category"] = filters.Category
	}

	if filters.Title != "" {
		filtersMap["title"] = filters.Title
	}

	if filters.Bound {
		filtersMap["lat"] = bson.M{
			"$gte": filters.MinLat,
			"$lte": filters.MaxLat,
		}

		filtersMap["long"] = bson.M{
			"$gte": filters.MinLong,
			"$lte": filters.MaxLong,
		}
	}

	if filters.Title != "" {
		filtersMap["title"] = filters.Title
	}

	filtersMap["hidden"] = filters.Hidden

	return filtersMap
}
