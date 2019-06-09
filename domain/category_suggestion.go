package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type CategorySuggestion struct {
	Id             primitive.ObjectID `json:"id" bson:"_id" required:"false"`
	CategoryName   string             `json:"name" bson:"name" required:"false"`
	HiddenCategory bool               `json:"hidden" bson:"hidden" required:"false"`
	CategoryIcon   string             `json:"iconClass" bson:"icon" required:"false"`
	Status         Status             `json:"status" bson:"status" required:"false"`
}

type Status string

const (
	Pending  Status = "waiting_for_approval"
	Approved Status = "approved"
	Rejected Status = "rejected"
)
