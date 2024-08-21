package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	Id        		primitive.ObjectID      `json:"_id" bson:"_id"`
	CustomerId        	  primitive.ObjectID      `json:"customerId" bson:"customerId"`
	CartItems		  []CartItem	         `json:"cartItems" bson:"cartItems"`
	NumberOfProduct 	int					`json:"numberOfProduct" bson:"numberOfProduct"`
	TotalAmount			int					`json:"totalAmount" bson:"totalAmount"`
}