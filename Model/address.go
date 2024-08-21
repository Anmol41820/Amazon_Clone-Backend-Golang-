package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

// import (
// 	"encoding/json"
// )

type Address struct {
	Id        			  primitive.ObjectID       `json:"_id" bson:"_id"`
	FullName        	string      `json:"fullName" bson:"fullName"`
	MobileNumber 		  string 		  `json:"mobileNumber" bson:"mobileNumber"`
	Pincode			 	    string		  `json:"pincode" bson:"pincode"`
	Line1        		  string      `json:"line1" bson:"line1"`
	Line2        		  string      `json:"line2" bson:"line2"`
	Landmark 			    string 		  `json:"landmark" bson:"landmark"`
	City			 	      string		  `json:"city" bson:"city"`
	State        		  string      `json:"state" bson:"state"`
	Country        		string      `json:"country" bson:"country"`
	IsDefault			bool			`json:"isDefault" bson:"isDefault"`
}

// type MyAddress struct{
// 	A Address
// }

// func (ass MyAddress) MarshalJSON() ([]byte, error) {
// 	m,_ := json.Marshal(ass.A)

// 	var a interface{}
// 	json.Unmarshal(m,&a)
// 	b := a.(map[string]interface{})

// 	b["id"] = b["_id"]
// 	delete(b, "_id")

// 	return json.Marshal(b)
// }