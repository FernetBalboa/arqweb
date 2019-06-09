package controller

import (
	"github.com/fernetbalboa/arqweb/src/apierror"
	"github.com/fernetbalboa/arqweb/src/domain"
	"github.com/fernetbalboa/arqweb/src/storage"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

const DefaultSearchLimit int64 = 20

type POIController struct {
	POIStorage storage.POIStorage
	CategoryStorage   storage.CategoryStorage
}

func CreatePOIController(POIStorage storage.POIStorage, CategoryStorage storage.CategoryStorage) *POIController {
	return &POIController{
		POIStorage: POIStorage,
		CategoryStorage:   CategoryStorage,
	}
}

func NewPOIController(categoryStorage storage.CategoryStorage) (*POIController, error) {
	POIStorage, err := storage.NewPOIStorage()
	if err != nil {
		return nil, apierror.Wrapf(err, "Could not create POI controller")
	}

	return CreatePOIController(POIStorage, categoryStorage), nil
}

func (pc *POIController) AddPOI(c *gin.Context) {
	var geoJsonFeature geojson.Feature
	err := c.ShouldBindJSON(&geoJsonFeature)

	if err != nil || geoJsonFeature.Geometry == nil {
		var apiError error
		errorMsg := "Error parsing POI. It should be a GeoJson feature"
		if err != nil {
			apiError = apierror.BadRequest.Wrapf(err, errorMsg)
		} else {
			apiError = apierror.BadRequest.Newf(errorMsg)
		}
		_ = c.Error(apiError)
		return
	}

	log.Infof("Adding new POI: %+v", geoJsonFeature)

	savedPOI, err := pc.POIStorage.SaveFeature(&geoJsonFeature)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, savedPOI)
}

func (pc *POIController) EditPOI(c *gin.Context) {
	id := c.Param("id")
	poiId, err := primitive.ObjectIDFromHex(id)

	var poi domain.PointOfInterest
	err = c.ShouldBindJSON(&poi)

	if err != nil {
		_ = c.Error(err)
		return
	}

	err = pc.POIStorage.EditPOI(poiId, &poi)

	c.JSON(http.StatusOK, gin.H{"id": poiId, "message": "POI successfully updated",
		"status": "OK", "code": http.StatusOK})
}

func (pc *POIController) SearchPOI(c *gin.Context) {
	var searchFilters domain.POIFilter

	if err := c.ShouldBindQuery(&searchFilters); err != nil {
		err = apierror.BadRequest.Wrapf(err, "Invalid search query filters")
		_ = c.Error(err)
		return
	}

	if searchFilters.Limit == 0 {
		searchFilters.Limit = DefaultSearchLimit
	}

	log.Infof("Searching POIs for request: %s", c.Request.URL)

	POIs, err := pc.POIStorage.SearchPOI(&searchFilters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var result []*domain.PointOfInterestDTO
	for _, poi := range POIs {
		filters := &domain.CategoryFilter{Name:poi.Category}
		category, err := pc.CategoryStorage.SearchCategory(filters)
		if err != nil {
			_ = c.Error(err)
			return
		}
		poiDTO := &domain.PointOfInterestDTO{
			Id: poi.Id,
			Title: poi.Title,
			Category: *category[0],
			Picture: poi.Picture,
			OwnerId: poi.OwnerId,
			Description: poi.Description,
			Lat: poi.Lat,
			Long: poi.Long,
			Hidden: poi.Hidden,
		}
		result = append(result, poiDTO)
	}

	log.Infof("Search results: %v", result)
	c.JSON(http.StatusOK, result)
}
