package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type DeliverOrder struct {
	Id        			 primitive.ObjectID       `json:"_id" bson:"_id"`
	CustomerId        	primitive.ObjectID      `json:"customerId" bson:"customerId"`
	ProductId        	  primitive.ObjectID      `json:"productId" bson:"productId"`
	OrderId				primitive.ObjectID      `json:"orderId" bson:"orderId"`
	CustomerName			string				`json:"customerName" bson:"customerName"`
	CustomerAddress			Address				`json:"customerAddress" bson:"customerAddress"`
	ExpectedDeliveryDate	time.Time			`json:"expectedDeliveryDate" bson:"expectedDeliveryDate"`
}