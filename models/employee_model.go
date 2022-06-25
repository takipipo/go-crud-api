package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id        primitive.ObjectID `json:"id"`
	Firstname string             `json:"firstname" validate:"required"`
	Lastname  string             `json:"lastname" validate:"required"`
	Position  string             `json:"position" validate:"required"`
	Age       int                `json:"age" validate:"required"`
}
