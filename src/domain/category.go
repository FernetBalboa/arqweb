package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id     primitive.ObjectID `json:"id" bson:"_id" required:"false"`
	Name   string             `json:"name" bson:"name" required:"false"`
	Hidden bool               `json:"hidden" bson:"hidden" required:"false"`
	Icon   string             `json:"iconClass" bson:"icon" required:"false"`
}
