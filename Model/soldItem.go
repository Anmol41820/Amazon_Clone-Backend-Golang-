package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SoldItem struct {
	ProductId        	    primitive.ObjectID      `json:"productId" bson:"productId"`
	ProductName				string				`json:"productName" bson:"productName"`
	ProductImage			[]string 				`json:"productImage" bson:"productImage"`
	ProductPrice			[]int					`json:"productPrice" bson:"productPrice"`
	Quantity      	        []int 		    `json:"quantity" bson:"quantity"`
	DeliveryDate			[]time.Time			`json:"deliveryDate" bson:"deliveryDate"`
}