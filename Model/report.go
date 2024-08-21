package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	Id        		primitive.ObjectID      `json:"_id" bson:"_id"`
	SellerId        	  primitive.ObjectID      `json:"sellerId" bson:"sellerId"`
	SoldItems		  []SoldItem	         `json:"soldItems" bson:"soldItems"`
	// NumberOfProduct 	int					`json:"numberOfProduct" bson:"numberOfProduct"`
}