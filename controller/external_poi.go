package controller

import (
	"github.com/fernetbalboa/arqweb/api/apierror"
	"github.com/fernetbalboa/arqweb/api/domain"
	"github.com/fernetbalboa/arqweb/api/storage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ExternalPOIController struct {
	ExternalPOIStorage storage.ExternalPOIStorage
	CategoryStorage   storage.CategoryStorage
}

func CreateExternalPOIController(ExternalPOIStorage storage.ExternalPOIStorage, CategoryStorage storage.CategoryStorage) *ExternalPOIController {
	return &ExternalPOIController{
		ExternalPOIStorage: ExternalPOIStorage,
		CategoryStorage:   CategoryStorage,
	}
}

func NewExternalPOIController(categoryStorage storage.CategoryStorage) (*ExternalPOIController, error) {
	ExtPOIStorage, err := storage.NewExternalPOIStorage()
	if err != nil {
		return nil, apierror.Wrapf(err, "Could not create External POI controller")
	}

	return CreateExternalPOIController(ExtPOIStorage, categoryStorage), nil
}

func (pc *ExternalPOIController) AddPOI(c *gin.Context) {
	var externalPoi domain.ExternalPOI
	err := c.ShouldBindJSON(&externalPoi)

	if err != nil || externalPoi.Id.IsZero() {
		var apiError error
		errorMsg := "Error parsing external POI."
		if err != nil {
			apiError = apierror.BadRequest.Wrapf(err, errorMsg)
		} else {
			apiError = apierror.BadRequest.Newf(errorMsg)
		}
		_ = c.Error(apiError)
		return
	}

	log.Infof("Adding new external POI: %+v", externalPoi)
	poi, err := pc.ExternalPOIStorage.SavePOI(&externalPoi)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, poi)
}

func (pc *ExternalPOIController) GetPOI(c *gin.Context) {
	id := c.Param("id")
	poiId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		_ = c.Error(err)
		return
	}

	poi, err := pc.ExternalPOIStorage.GetPOI(poiId)

	c.JSON(http.StatusOK, poi)
}


func (pc *ExternalPOIController) RemovePOI(c *gin.Context) {
	id := c.Param("id")
	poiId, err := primitive.ObjectIDFromHex(id)

	res, err := pc.ExternalPOIStorage.RemovePOI(poiId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (pc *ExternalPOIController) GetPOIs(c *gin.Context) {
	pois, err := pc.ExternalPOIStorage.GetPOIs()

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pois)
}

