package controller

import (
	"github.com/fernetbalboa/arqweb/src/api/apierror"
	"github.com/fernetbalboa/arqweb/src/api/domain"
	"github.com/fernetbalboa/arqweb/src/api/storage"
	"github.com/gin-gonic/gin"
	"github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const DefaultSearchLimit int64 = 20

type POIController struct {
	POIStorage storage.POIStorage
}

func CreatePOIController(POIStorage storage.POIStorage) *POIController {
	return &POIController{
		POIStorage: POIStorage,
	}
}

func NewPOIController() (*POIController, error) {
	POIStorage, err := storage.NewPOIStorage()
	if err != nil {
		return nil, apierror.Wrapf(err, "Could not create POI controller")
	}

	return CreatePOIController(POIStorage), nil
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
	poiId := c.Param("id")
	var poi domain.PointOfInterest
	err := c.ShouldBindJSON(&poi)

	if err != nil {
		_ = c.Error(err)
		return
	}

	if poiId != poi.Id.String() {
		apiError := apierror.BadRequest.New("POI ID can not be updated")
		_ = c.Error(apiError)
		return
	}

	err = pc.POIStorage.EditPOI(&poi)

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

	log.Infof("Search results: %v", POIs)
	c.JSON(http.StatusOK, POIs)
}

