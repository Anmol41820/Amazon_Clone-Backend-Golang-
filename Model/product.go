package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id        			    primitive.ObjectID      `json:"_id" bson:"_id"`
	ProductName			    string 		  `json:"productName" bson:"productName"`
	AboutProduct 		    string 		  `json:"aboutProduct" bson:"aboutProduct"`
	SellerId			    primitive.ObjectID		  `json:"sellerId" bson:"sellerId"`
	Brand				     string		  `json:"brand" bson:"brand"`
	Color					string			`json:"color" bson:"color"`
	ReleaseDate				time.Time         `json:"releaseDate" bson:"releaseDate"`
	BestSeller				bool			`json:"bestSeller" bson:"bestSeller"`
	NewRelease				bool 			`json:"newRelease" bson:"newRelease"`
	ReplacePolicy			bool			`json:"replacePolicy" bson:"replacePolicy"`
	ReturnPolicy			bool			`json:"returnPolicy" bson:"returnPolicy"`
	// Notify					[]primitive.ObjectID   `json:"notify" bson:"notify"`
	MaxRetailPrice		    int 		    `json:"maxRetailPrice" bson:"maxRetailPrice"`
	SellingPrice		    int			    `json:"sellingPrice" bson:"sellingPrice"`
	Discount				float64 			`json:"discount" bson:"discount"`
	Quantity			      int 		    `json:"quantity" bson:"quantity"`
	UnitsSold			      int 		    `json:"unitsSold" bson:"unitsSold"`
	ProductCategories	  []string 	  `json:"productCategories" bson:"productCategories"`
	AverageRating		    float64 		    `json:"averageRating" bson:"averageRating"`
	ProductImages 		  []string 	  `json:"productImages" bson:"productImages"`
	ProductProperties    map[string]string 	  `json:"productProperties" bson:"productProperties"`
	// Reviews				 []Review	  `json:"reviews" bson:"reviews"`
	NumberOfReviews			int		`json:"numberOfReviews" bson:"numberOfReviews"`
}