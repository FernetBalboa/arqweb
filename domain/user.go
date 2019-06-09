package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id" required:"false"`
	Username string             `json:"username" bson:"username" required:"true"`
	Password string             `json:"password" bson:"password" required:"true"`
	IsAdmin  bool               `json:"isAdmin" bson:"isAdmin" required:"false"`
}
