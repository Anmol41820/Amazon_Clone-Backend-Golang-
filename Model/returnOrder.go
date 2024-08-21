package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type ReturnOrder struct {
	Id        			 primitive.ObjectID       `json:"_id" bson:"_id"`
	CustomerId        	primitive.ObjectID      `json:"customerId" bson:"customerId"`
	ProductId        	  primitive.ObjectID      `json:"productId" bson:"productId"`
	OrderId				primitive.ObjectID      `json:"orderId" bson:"orderId"`
	CustomerName			string				`json:"customerName" bson:"customerName"`
	CustomerAddress			Address				`json:"customerAddress" bson:"customerAddress"`
	ExpectedReturnDate	   time.Time			`json:"expectedReturnDate" bson:"expectedReturnDate"`
	IsDamage			  bool			`json:"isDamage" bson:"isDamage"`
	DontLikeDueToColorOrSize		bool   `json:"dontLikeDueToColorOrSize" bson:"dontLikeDueToColorOrSize"`
}