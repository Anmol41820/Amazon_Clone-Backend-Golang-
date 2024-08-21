package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	ProductId        	    primitive.ObjectID      `json:"productId" bson:"productId"`
	ProductName				string				`json:"productName" bson:"productName"`
	ProductPrice			int					`json:"productPrice" bson:"productPrice"`
	QuantityInCart 		    int 		    `json:"quantityInCart" bson:"quantityInCart"`
	SelectedForBuying 	  bool		    `json:"selectedForBuying" bson:"selectedForBuying"`
	ProductImage			string 				`json:"productImage" bson:"productImage"`
	ProductInStock			bool				`json:"productInStock" bson:"productInStock"`
}