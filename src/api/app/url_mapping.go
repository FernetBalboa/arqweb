package app

import (
	"github.com/fernetbalboa/arqweb/src/api/controller"
	log "github.com/sirupsen/logrus"
)

var poiController *controller.POIController

func init() {
	POIController, err := controller.NewPOIController()
	if err != nil {
		log.Fatalf("Could not create POI controller. Cause: %s", err.Error())
	}
	poiController = POIController
}

// LoadEndpoints is the base function to map endpoints.
func LoadEndpoints() {
	Router.GET("/ping", controller.Ping)

	poiGroup := Router.Group("/poi")

	poiGroup.POST("", poiController.AddPOI)
	poiGroup.GET("/search", poiController.SearchPOI)
	poiGroup.PUT("/:id", poiController.EditPOI)

	categoryGroup := Router.Group("/category")
	categoryGroup.GET("", poiController.GetCategories)
	categoryGroup.POST("", poiController.AddCategory)
	categoryGroup.PUT("/:id", poiController.EditCategory)
}
