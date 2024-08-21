package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SearchHistory struct {
	Id        		primitive.ObjectID      `json:"_id" bson:"_id"`
	CustomerId        	  primitive.ObjectID      `json:"customerId" bson:"customerId"`
	SearchText			[]string					`json:"searchText" bson:"searchText"`
}