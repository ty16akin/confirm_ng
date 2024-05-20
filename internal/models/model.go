package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID   primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	FirtName string             `json:"first_name" bson:"first_name"`
	Lastname string             `json:"last_name" bson:"last_name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Status   string             `json:"status" bson:"status"`
	Created  *time.Time         `json:"created" bson:"created"`
}
