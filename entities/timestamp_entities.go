package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatedBy struct {
	UserId				primitive.ObjectID 		`json:"userId" bson:"userId"`
	Name 				string 					`json:"name" bson:"name"`
	Date 				time.Time				`json:"date" bson:"date"`
	Description 		string 					`json:"description" bson:"description"`
}

type UpdatedBy struct {
	UserId 				*primitive.ObjectID		`json:"userId" bson:"userId"`
	Name	 			*string 				`json:"name" bson:"name"`
	Date 				*time.Time				`json:"date" bson:"date"`
	Description 		*string 				`json:"description" bson:"description"`
}

type DeletedBy struct {
	UserId 				primitive.ObjectID 	`json:"userId" bson:"userId"`
	Name	 			string 				`json:"name" bson:"name"`
	Date 				time.Time			`json:"date" bson:"date"`
	Description 		string 				`json:"description" bson:"description"`
}
