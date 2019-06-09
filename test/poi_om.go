package test

import (
	"github.com/fernetbalboa/arqweb/domain"
	"github.com/paulmach/go.geojson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DefaultGeoJsonFeatureString() string {
	return `{
		"type": "Feature",
			"geometry": {
			"type": "Point",
				"coordinates": [20, 40]
	},
	"properties": {
	"title": "testTitle",
	"category": "testCategory",
	"type": "testType",
	"description": "testDescription",
	"hidden": false
	}
}`
}

func DefaultGeoJsonFeature() *geojson.Feature {
	geoJsonPoint := geojson.NewFeature(geojson.NewPointGeometry([]float64{20, 40}))
	geoJsonPoint.SetProperty("title", "testTitle")
	geoJsonPoint.SetProperty("category", "testCategory")
	geoJsonPoint.SetProperty("description", "testDescription")
	geoJsonPoint.SetProperty("type", "testType")
	geoJsonPoint.SetProperty("hidden", false)

	return geoJsonPoint
}

func DefaultPOI() *domain.PointOfInterest {
	defGeoJson := DefaultGeoJsonFeature()
	return &domain.PointOfInterest{
		Id:          NewDocumentId(),
		Title:       defGeoJson.PropertyMustString("title", ""),
		Category:    defGeoJson.PropertyMustString("category", ""),
		Description: defGeoJson.PropertyMustString("description", ""),
		Lat:         defGeoJson.Geometry.Point[0],
		Long:        defGeoJson.Geometry.Point[1],
		Hidden:      false,
	}
}

func NewDocumentId() primitive.ObjectID {
	return primitive.NewObjectID()
}
