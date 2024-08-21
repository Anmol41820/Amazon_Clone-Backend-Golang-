package address

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"

	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAddress(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var address []Model.Address

	//finding address in database
	collection := conn.ConnectDB("addresses")
	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Address
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		address = append(address, u)
	}
	json.NewEncoder(w).Encode(address)
}

func AddAddress(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var params = mux.Vars(r)
		ids := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		address := make(map[string]interface{})
		address["isDefault"] = false
		err = json.Unmarshal([]byte(asString), &address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if address == nil {
			http.Error(w, "Invalid address data", http.StatusBadRequest)
			return
		}
		// delete(address, "_id")
		newAddressId := primitive.NewObjectID()
		address["_id"] = newAddressId

		//validation for all the fields
		if address["fullName"].(string) == ""{
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintln(w,"FullName is empty!!")
			return
		}
		mobileNumber, ok := address["mobileNumber"].(string)
		if !ok {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if !validationForMobileNumber(mobileNumber, w) {
			return
		}
		if address["pincode"].(string) == ""{
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintln(w,"Pincode is empty!!")
			return
		}
		if address["line1"].(string) == ""{
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintln(w,"Line1 is empty!!")
			return
		}
		if address["landmark"].(string) == ""{
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintln(w,"Landmark is empty!!")
			return
		}
		if address["state"].(string) == ""{
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintln(w,"State is empty!!")
			return
		}
		if address["country"].(string) == ""{
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintln(w,"Country is empty!!")
			return
		}

		//connection to addressess DB and inserting the address in it
		collection := conn.ConnectDB("addresses")
		result, err := collection.InsertOne(context.TODO(), address)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//changing the address format, to add it in the customer
		var newAddress Model.Address
		newAddress.Id = newAddressId
		newAddress.FullName = address["fullName"].(string)
		newAddress.City = address["city"].(string)
		newAddress.Country = address["country"].(string)
		newAddress.IsDefault = address["isDefault"].(bool)
		newAddress.Landmark = address["landmark"].(string)
		newAddress.Line1 = address["line1"].(string)
		newAddress.Line2 = address["line2"].(string)
		newAddress.MobileNumber = address["mobileNumber"].(string)
		newAddress.Pincode = address["pincode"].(string)
		newAddress.State = address["state"].(string)


		//finding the customer and updating/adding the address
		var customer Model.Customer
		customerCollection := conn.ConnectDB("customers")
		filter := bson.M{"_id": customerId}
		errr := customerCollection.FindOne(context.TODO(), filter).Decode(&customer)
		if errr != nil {
			conn.GetError(err, w)
			return
		}
		if newAddress.IsDefault{
			for i:=0;i<len(customer.Addresses);i++{
				customer.Addresses[i].IsDefault = false
			}
		}
		customer.Addresses = append(customer.Addresses, newAddress)
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "addresses", Value: customer.Addresses},
				},
			},
		}
		customerResult, errrr := customerCollection.UpdateOne(context.TODO(), filter, update)
		if errrr != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(customerResult)
	}
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	// if(!conn.ProtectedHandler(w,r)){return}

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["addressId"]
		id, _ := primitive.ObjectIDFromHex(ids)
		var address Model.Address
		address.Id = id
		filter := bson.M{"_id": id}

		_ = json.NewDecoder(r.Body).Decode(&address)

		//validation for all the fields
		if address.FullName == ""{
			fmt.Fprintln(w,"FullName is empty!!")
			return
		}
		mobileNumber:= address.MobileNumber
		if !validationForMobileNumber(mobileNumber, w) {
			return
		}
		if address.Pincode == ""{
			fmt.Fprintln(w,"Pincode is empty!!")
			return
		}
		if address.Line1 == ""{
			fmt.Fprintln(w,"Line1 is empty!!")
			return
		}
		if address.Landmark == ""{
			fmt.Fprintln(w,"Landmark is empty!!")
			return
		}
		if address.State == ""{
			fmt.Fprintln(w,"State is empty!!")
			return
		}
		if address.Country == ""{
			fmt.Fprintln(w,"Country is empty!!")
			return
		}

		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "fullName", Value: address.FullName},
					{Key: "mobileNumber", Value: address.MobileNumber},
					{Key: "pincode", Value: address.Pincode},
					{Key: "line1", Value: address.Line1},
					{Key: "line2", Value: address.Line2},
					{Key: "landmark", Value: address.Landmark},
					{Key: "city", Value: address.City},
					{Key: "state", Value: address.State},
					{Key: "country", Value: address.Country},
					{Key: "isDefault", Value: address.IsDefault},
				},
			},
		}

		collection := conn.ConnectDB("addresses")
		result,err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//updating the address in the customer model
		idd := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(idd)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var customer Model.Customer
		customerCollection := conn.ConnectDB("customers")
		customerFilter := bson.M{"_id": customerId}
		errr := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
		if errr != nil {
			conn.GetError(err, w)
			return
		}
		for i:=0;i<len(customer.Addresses);i++{
			if customer.Addresses[i].Id == address.Id{
				customer.Addresses[i] = address
				customer.Addresses[i].IsDefault = true
			}else{
				customer.Addresses[i].IsDefault = false
			}
		}
		customerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "addresses", Value: customer.Addresses},
				},
			},
		}
		customerResult, errrr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
		if errrr != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(customerResult)
	}
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["addressId"]
		id, _ := primitive.ObjectIDFromHex(ids)

		filter := bson.M{"_id": id}
		collection := conn.ConnectDB("addresses")

		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//deleting the address in the customer model
		idd := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(idd)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var customer Model.Customer
		customerCollection := conn.ConnectDB("customers")
		customerFilter := bson.M{"_id": customerId}
		errr := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
		if errr != nil {
			conn.GetError(err, w)
			return
		}
		for i:=0;i<len(customer.Addresses);i++{
			if customer.Addresses[i].Id == id{
				customer.Addresses = append(customer.Addresses[:i], customer.Addresses[i+1:]...)
				break
			}
		}
		for i:=0;i<len(customer.Addresses);i++{
			customer.Addresses[i].IsDefault = false
		}
		if len(customer.Addresses) > 0{
			customer.Addresses[0].IsDefault = true
		}
		customerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "addresses", Value: customer.Addresses},
				},
			},
		}
		customerResult, errrr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
		if errrr != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(deleteResult)
		json.NewEncoder(w).Encode(customerResult)
	}

}


//validation for mobile number
func validationForMobileNumber(mobileNumber string, w http.ResponseWriter) bool {
	if len(mobileNumber) == 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("Mobile number is Empty!!")
		return false
	}
	if len(mobileNumber) != 10 {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Invalid Mobile number!!")
		json.NewEncoder(w).Encode("Your Mobile number must be of 10 digits!!")
		return false
	}
	for i := 0; i < len(mobileNumber); i++ {
		if mobileNumber[i] < '0' || mobileNumber[i] > '9' {
			w.WriteHeader(http.StatusConflict)
			fmt.Println("Invalid Mobile number, mobile number should not be character!!")
			json.NewEncoder(w).Encode("Invalid Mobile number, mobile number should not be character!!")
			return false
		}
	}
	return true
}