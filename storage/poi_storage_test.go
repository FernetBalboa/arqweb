package storage_test

import (
	"context"
	"github.com/fernetbalboa/arqweb/config"
	"github.com/fernetbalboa/arqweb/domain"
	"github.com/fernetbalboa/arqweb/mock"
	"github.com/fernetbalboa/arqweb/storage"
	"github.com/fernetbalboa/arqweb/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const (
	TestDB                 string = "test_db"
	POITestCollection             = "poi_test"
	CategoryTestCollection        = "cat_test"
)

func init() {
	config.Scope = config.TestScope
}

func TestSaveGeoJsonFeature(t *testing.T) {
	//Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	poiCollectionMock := mock.NewMockICollection(ctrl)

	POIStorage, _ := storage.CreatePOIStorage(poiCollectionMock)

	pointFeature := test.DefaultGeoJsonFeature()
	documentId := test.NewDocumentId()

	poiCollectionMock.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(&mongo.InsertOneResult{
		InsertedID: documentId,
	}, nil)

	//When
	savedPOI, err := POIStorage.SaveFeature(pointFeature)

	//Then
	assert.NoError(t, err)
	assert.Equal(t, documentId, savedPOI.Id)
	assert.Equal(t, pointFeature.PropertyMustString("title"), savedPOI.Title)
	assert.Equal(t, pointFeature.PropertyMustString("category"), savedPOI.Category)
	assert.Equal(t, pointFeature.PropertyMustString("description"), savedPOI.Description)
	assert.Equal(t, pointFeature.Geometry.Point[0], savedPOI.Lat)
	assert.Equal(t, pointFeature.Geometry.Point[1], savedPOI.Long)
}

//Integration test using local MongoDB
func TestSearchPOI(t *testing.T) {
	//Given
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	POICollection := client.Database(TestDB).Collection(POITestCollection)

	POIStorage, _ := storage.CreatePOIStorage(POICollection)

	t.Run("Search POI by category", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database(TestDB).Drop(context.Background())
		POI := test.DefaultPOI()
		anotherPOI := test.DefaultPOI()
		anotherPOI.Category = "anotherCategory"

		savedPOI, saveErr1 := POIStorage.SavePOI(POI)
		anotherCategoryPOI, saveErr2 := POIStorage.SavePOI(anotherPOI)

		categoryFilter := &domain.POIFilter{Category: POI.Category}

		//When
		foundPOIs, searchErr := POIStorage.SearchPOI(categoryFilter)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundPOIs, savedPOI)
		assert.NotContains(t, foundPOIs, anotherCategoryPOI)
	})

	t.Run("Search POI by hidden value", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database(TestDB).Drop(context.Background())
		POI := test.DefaultPOI()
		anotherPOI := test.DefaultPOI()
		anotherPOI.Hidden = true

		savedPOI, saveErr1 := POIStorage.SavePOI(POI)
		hiddenPoi, saveErr2 := POIStorage.SavePOI(anotherPOI)

		hiddenFilter := &domain.POIFilter{Hidden: POI.Hidden}

		//When
		foundPOIs, searchErr := POIStorage.SearchPOI(hiddenFilter)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundPOIs, savedPOI)
		assert.NotContains(t, foundPOIs, hiddenPoi)
	})

	t.Run("Search POI in latitude bounds", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database(TestDB).Drop(context.Background())
		insideBoundsPOI := test.DefaultPOI()
		minLat := insideBoundsPOI.Lat - 10
		maxLat := insideBoundsPOI.Lat + 10
		minLong := insideBoundsPOI.Long - 10
		maxLong := insideBoundsPOI.Long + 10

		outsideBoundsPOI := test.DefaultPOI()
		outsideBoundsPOI.Lat = minLat - 1

		savedInsideBoundsPOI, saveErr1 := POIStorage.SavePOI(insideBoundsPOI)
		savedOutsideBoundsPOI, saveErr2 := POIStorage.SavePOI(outsideBoundsPOI)

		boundsFilter := &domain.POIFilter{
			MaxLat:  maxLat,
			MaxLong: maxLong,
			MinLat:  minLat,
			MinLong: minLong,
			Bound:   true,
		}

		//When
		foundPOIs, searchErr := POIStorage.SearchPOI(boundsFilter)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundPOIs, savedInsideBoundsPOI)
		assert.NotContains(t, foundPOIs, savedOutsideBoundsPOI)
	})

	t.Run("Search POI in longitude bounds", func(t *testing.T) {
		//noinspection GoUnhandledErrorResult
		defer client.Database(TestDB).Drop(context.Background())
		insideBoundsPOI := test.DefaultPOI()
		minLat := insideBoundsPOI.Lat - 10
		maxLat := insideBoundsPOI.Lat + 10
		minLong := insideBoundsPOI.Long - 10
		maxLong := insideBoundsPOI.Long + 10

		outsideBoundsPOI := test.DefaultPOI()
		outsideBoundsPOI.Long = maxLong + 1

		savedInsideBoundsPOI, saveErr1 := POIStorage.SavePOI(insideBoundsPOI)
		savedOutsideBoundsPOI, saveErr2 := POIStorage.SavePOI(outsideBoundsPOI)

		boundsFilter := &domain.POIFilter{
			MaxLat:  maxLat,
			MaxLong: maxLong,
			MinLat:  minLat,
			MinLong: minLong,
			Bound:   true,
		}

		//When
		foundPOIs, searchErr := POIStorage.SearchPOI(boundsFilter)

		//Then
		assert.NoError(t, saveErr1)
		assert.NoError(t, saveErr2)
		assert.NoError(t, searchErr)

		assert.Contains(t, foundPOIs, savedInsideBoundsPOI)
		assert.NotContains(t, foundPOIs, savedOutsideBoundsPOI)
	})
}
