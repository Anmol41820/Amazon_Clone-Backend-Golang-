package Model

type Wishlist struct {
	Id        		  string     `json:"_id" bson:"_id"`
	CustomerId        	  string      `json:"customerId" bson:"customerId"`
	WishlistItems 		[]WishlistItem 		    `json:"wishlistItems" bson:"wishlistItems"`
	NumberOfProduct 	int					`json:"numberOfProduct" bson:"numberOfProduct"`
}