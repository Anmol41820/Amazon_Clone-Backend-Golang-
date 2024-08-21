package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id        		  primitive.ObjectID      				      `json:"_id" bson:"_id"`
	CategoryName	  string 						        `json:"categoryName" bson:"categoryName"`
	Brands        	[]string    				      `json:"brands" bson:"brands"`
	Colors        	[]string    				      `json:"colors" bson:"colors"`	
	PriceRanges 	  map[string]string					        `json:"priceRanges" bson:"priceRanges"`
	// SellerIds 		  []string					        `json:"sellerIds" bson:"sellerIds"`
	// Properties 		  map[string]string 			`json:"properties" bson:"properties"`
}