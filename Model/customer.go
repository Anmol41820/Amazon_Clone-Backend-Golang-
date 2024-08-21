package Model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	Id        						        primitive.ObjectID       `json:"_id" bson:"_id"`
	FirstName 						        string   	  `json:"firstName" bson:"firstName" `
	LastName  						        string   	  `json:"lastName" bson:"lastName" `
	Email     						        string   	  `json:"email" bson:"email" `
	Password  						        string   	  `json:"password" bson:"password" `
	MobileNumber    				      string   	  `json:"mobileNumber" bson:"mobileNumber"`
	DateOfBirth     				      string   	  `json:"dateOfBirth" bson:"dateOfBirth"`
	Role      						        string   	  `json:"role" bson:"role"` 
	IsPrime									bool		`json:"isPrime" bson:"isPrime"`
	PrimeExpireDate						time.Time		`json:"primeExpireDate" bson:"primeExpireDate"`
	ProductsPurchased 				    []Order 	  `json:"productsPurchased" bson:"productsPurchased"`
	Addresses 		 				        []Address 	`json:"addresses" bson:"addresses"`
	Wallet								int 		`json:"wallet" bson:"wallet"`
	// FriendRequests 					      []string 	  `json:"friendRequests" bson:"friendRequests"`
	// Friends 						          []string 	  `json:"friends" bson:"friends"`
	// FriendProductRecommendations 	[]string 	  `json:"friendProductRecommendations" bson:"friendProductRecommendations"`
}