package Model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type PriceDetail struct {
	ListPrice           int   `json:"listPrice" bson:"listPrice"`
	SellingPrice		int   `json:"sellingPrice" bson:"sellingPrice"`
	DeliveryCharge		int   `json:"deliveryCharge" bson:"deliveryCharge"`
	TotalAmount			int 	`json:"totalAmount" bson:"totalAmount"`
}

type Order struct {
	Id        			 primitive.ObjectID       `json:"_id" bson:"_id"`
	CustomerId        	primitive.ObjectID      `json:"customerId" bson:"customerId"`
	ProductIds        	  []primitive.ObjectID      `json:"productIds" bson:"productIds"`
	ProductNames			[]string			`json:"productNames" bson:"productNames"`
	OrderQuantitys		    []int			    `json:"orderQuantitys" bson:"orderQuantitys"`
	OrderedDate			    time.Time   	  `json:"orderedDate" bson:"orderedDate"`
	DeliveredDates		    [][]time.Time   	  `json:"deliveredDates" bson:"deliveredDates"`
	Status		    	    [][]string   	  `json:"status" bson:"status"`
	ShippingAddress    	Address   	`json:"shippingAddress" bson:"shippingAddress"`
	PaymentDetails		  Payment   	`json:"paymentDetails" bson:"paymentDetails"` 
	PriceDetail				PriceDetail    `json:"priceDetail" bson:"priceDetail"` 
	TotalAmount 		    int 		    `json:"totalAmount" bson:"totalAmount"`
}