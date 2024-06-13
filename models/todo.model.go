package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID         primitive.ObjectID `bson:"_id"`
	Todo       string             `json:"todo" validate:"required"`
	Todo_id    string             `bson:"todo_id"`
	Completed  *bool              `json:"completed"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
