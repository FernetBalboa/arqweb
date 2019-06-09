package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PointOfInterest struct {
	Id          primitive.ObjectID `json:"id" bson:"_id" required:"false"`
	Title       string             `json:"title" bson:"title" required:"true"`
	Category    string             `json:"category" bson:"category" required:"true"`
	Picture     Picture            `json:"picture" bson:"picture" required:"false"`
	OwnerId     int64              `json:"ownerId" bson:"ownerId" required:"false"`
	Description string             `json:"description" bson:"description" required:"false"`
	Lat         float64            `json:"lat" bson:"lat" required:"false"`
	Long        float64            `json:"long" bson:"long" required:"false"`
	Hidden      bool               `json:"hidden" bson:"hidden" required:"false"`
}

type Picture struct {
	Name        string `json:"name" bson:"name" required:"false"`
	Data        string `json:"data" bson:"data" required:"false"`
	ContentType string `json:"contentType" bson:"contentType" required:"false"`
}

type PointOfInterestDTO struct {
	Id          primitive.ObjectID `json:"id" bson:"_id" required:"false"`
	Title       string             `json:"title" bson:"title" required:"true"`
	Category    Category           `json:"category" bson:"category" required:"true"`
	Picture     Picture            `json:"picture" bson:"picture" required:"false"`
	OwnerId     int64              `json:"ownerId" bson:"ownerId" required:"false"`
	Description string             `json:"description" bson:"description" required:"false"`
	Lat         float64            `json:"lat" bson:"lat" required:"false"`
	Long        float64            `json:"long" bson:"long" required:"false"`
	Hidden      bool               `json:"hidden" bson:"hidden" required:"false"`
}