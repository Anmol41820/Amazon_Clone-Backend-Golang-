package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type WishlistItem struct {
	ProductId        	    primitive.ObjectID      `json:"productId" bson:"productId"`
	ProductName				string				`json:"productName" bson:"productName"`
	ProductImage			string 				`json:"productImage" bson:"productImage"`
	ProductPrice			int					`json:"productPrice" bson:"productPrice"`
	ProductInStock			bool				`json:"productInStock" bson:"productInStock"`
}