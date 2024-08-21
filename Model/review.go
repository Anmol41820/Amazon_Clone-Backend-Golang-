package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	Id        			    primitive.ObjectID      `json:"_id" bson:"_id"`
	ProductId        	  primitive.ObjectID      `json:"productId" bson:"productId"`
	CustomerId        	primitive.ObjectID      `json:"customerId" bson:"customerId"`
	CustomerName		string					`json:"customerName" bson:"customerName"`
	Rating 				      int 			  `json:"rating" bson:"rating"`
	Headline 			      string			`json:"headline" bson:"headline"`
	Description 		    string			`json:"description" bson:"description"`
	ReviewDate 				   time.Time		`json:"reviewDate" bson:"reviewDate"`
	ReviewImages 		    []string		`json:"reviewImages" bson:"reviewImages"`
}