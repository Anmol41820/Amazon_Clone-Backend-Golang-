package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CustomerCare struct {
	Id        						        primitive.ObjectID       `json:"_id" bson:"_id"`
	CustomerId								primitive.ObjectID		`json:"customerId" bson:"customerId"`
	MessageFromCustomer 					[]string   	  `json:"messageFromCustomer" bson:"messageFromCustomer" `
	MessageFromCustomerCare  				[]string   	  `json:"messageFromCustomerCare" bson:"messageFromCustomerCare" `
}