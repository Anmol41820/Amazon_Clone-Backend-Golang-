package home

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	"encoding/json"
	"fmt"
	"time"

	"Amazon_Server/Model"
	"context"

	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckPrimeMemberShip(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")

	//updating the customer prime membership
	var params = mux.Vars(r)
	ids := params["id"]
	customerId, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	customerCollection := conn.ConnectDB("customers")
	customerFilter := bson.M{"_id": customerId}
	var customer Model.Customer
	errr := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}

	if customer.PrimeExpireDate.Before(time.Now()){
		customer.IsPrime = false
		customerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "isPrime", Value: customer.IsPrime},
				},
			},
		}
	
		customerResult,err := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		if customer.IsPrime{
			fmt.Fprintln(w,"Prime Membership : True")
		}else{
			fmt.Fprintln(w,"Prime Membership : False")
		}
		json.NewEncoder(w).Encode(customerResult)
	}
	
}
