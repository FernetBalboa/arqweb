package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ExternalPOI struct {
	Id      primitive.ObjectID `json:"id" bson:"_id"`
	Owner 	string `json:"owner" bson:"owner"`
}
