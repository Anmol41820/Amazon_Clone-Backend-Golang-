package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	Id        			 primitive.ObjectID      `json:"_id" bson:"_id"`
	PaymentMethod       string      `json:"paymentMethod" bson:"paymentMethod"`
	TotalAmount 		    int 		    `json:"totalAmount" bson:"totalAmount"`	
	UserId        		  primitive.ObjectID      `json:"userId" bson:"userId"`
	CardNumber			    string   	  `json:"cardNumber" bson:"cardNumber"`
	CardExpiryDate		  string 		  `json:"cardExpiryDate" bson:"cardExpiryDate"`
	NameOnCard			   string   	  `json:"nameOnCard" bson:"nameOnCard"`
	UpiId		    	    string   	  `json:"upiId" bson:"upiId"`
}