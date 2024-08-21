package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Seller struct {
	Id        		        primitive.ObjectID     `json:"_id" bson:"_id"`
	FirstName 		        string     `json:"firstName" bson:"firstName" `
	LastName  		        string     `json:"lastName" bson:"lastName" `
	Email     		        string     `json:"email" bson:"email" `
	Password  		        string     `json:"password" bson:"password" `
	MobileNumber          string     `json:"mobileNumber" bson:"mobileNumber"`
	DateOfBirth           string     `json:"dateOfBirth" bson:"dateOfBirth"`
	Role      		        string     `json:"role" bson:"role"`
	IsPrime    				bool 		`json:"isPrime" bson:"isPrime"`
	ProductsListedIds     []primitive.ObjectID   `json:"productsListedIds" bson:"productsListedIds"`
	ProductsSoldIds       []primitive.ObjectID   `json:"productsSoldIds" bson:"productsSoldIds"`
}
