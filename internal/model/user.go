package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Created  time.Time          `bson:"created" json:"created"`
	Updated  time.Time          `bson:"updated" json:"updated"`
}

type CreateUserRequest struct {
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
