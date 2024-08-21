package Model

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecentlyViewedProduct struct {
	Id        			    primitive.ObjectID      `json:"_id" bson:"_id"`
	CustomerId				primitive.ObjectID		`json:"customerId" bson:"customerId"`
	// Products		    []Product 		  `json:"products" bson:"products"`
	ProductIds				[]primitive.ObjectID	`json:"productIds" bson:"productIds"`
	ProductName			    []string 		  `json:"productName" bson:"productName"`
	AboutProduct 		    []string 		  `json:"aboutProduct" bson:"aboutProduct"`
	Brand				    []string		  `json:"brand" bson:"brand"`
	BestSeller				[]bool			`json:"bestSeller" bson:"bestSeller"`
	NewRelease				[]bool 			`json:"newRelease" bson:"newRelease"`
	MaxRetailPrice		    []int 		    `json:"maxRetailPrice" bson:"maxRetailPrice"`
	SellingPrice		    []int			    `json:"sellingPrice" bson:"sellingPrice"`
	Discount				[]float64 			`json:"discount" bson:"discount"`
	AverageRating		    []float64 		    `json:"averageRating" bson:"averageRating"`
	ProductImage		  	[]string 	  	`json:"productImage" bson:"productImage"`
	NumberOfReviews			[]int			`json:"numberOfReviews" bson:"numberOfReviews"`
}