package seller

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	helper "Amazon_Server/Helper"
	"time"

	"Amazon_Server/Model"
	"context"
	"fmt"

	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSeller(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var seller []Model.Seller

	collection := conn.ConnectDB("sellers")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Seller
		err := curr.Decode(&u)
		if err != nil {
			// log.Fatal(err)
			w.WriteHeader(http.StatusBadGateway)
		}
		if err != nil {
			fmt.Println("****ERROR*****")
			// log.Fatal(err)
			w.WriteHeader(http.StatusBadGateway)
		}
		seller = append(seller, u)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(seller)
}

func CreateSeller(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		seller := make(map[string]interface{})
		seller["role"] = "seller"
		seller["isPrime"] = false

		err = json.Unmarshal([]byte(asString), &seller)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if seller == nil {
			http.Error(w, "Invalid seller data", http.StatusBadRequest)
			return
		}
		// delete(seller, "_id")

		//check wheather the register email or mobile number is already exist or not
		collection := conn.ConnectDB("sellers")

		coll, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer coll.Close(context.TODO())

		for coll.Next(context.TODO()) {
			var existingseller Model.Seller
			err := coll.Decode(&existingseller)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if existingseller.Id == seller["_id"] {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Duplicate Id!!")
				return
			}
			if existingseller.Email == seller["email"].(string) || existingseller.MobileNumber == seller["mobileNumber"].(string) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Email or Mobile number already used!!")
				return
			}
		}
		if err := coll.Err(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		//Validation for all field like email, password, mobilenumber etc..
		if !validationForRegister(seller, w) {
			return
		}

		//register -> encry the password
		password, ok := seller["password"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encryptedValue := helper.Encrypt([]byte(password), "Secret Key")
		str_encryptedVal := hex.EncodeToString(encryptedValue)
		seller["password"] = str_encryptedVal
		seller["_id"] = primitive.NewObjectID()

		// collection := conn.ConnectDB("sellers")
		result, err := collection.InsertOne(context.TODO(), seller)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		userCollection := conn.ConnectDB("users")
		userResult, err := userCollection.InsertOne(context.TODO(), seller)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//creating a new report for the new seller
		newReport := make(map[string]interface{})
		newReport["_id"] = primitive.NewObjectID()
		newReport["sellerId"] = seller["_id"]
		newReport["soldItems"] = []Model.SoldItem{}
		reportCollection := conn.ConnectDB("reports")
		reportResult, err := reportCollection.InsertOne(context.TODO(), newReport)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(userResult)
		json.NewEncoder(w).Encode(reportResult)
	}
}

func GetSingleSeller(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("sellers")

	var seller Model.Seller
	var params = mux.Vars(r)

	ids := params["id"]
	id, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&seller)

	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(seller)
}

func UpdateSeller(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	// if(!conn.ProtectedHandler(w,r)){return}

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["id"]
		id, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var seller Model.Seller

		filter := bson.M{"_id": id}

		_ = json.NewDecoder(r.Body).Decode(&seller)
		// seller.IsPrime = false
		// seller.Addresses = []Model.Address{}
		// seller.ProductsPurchased = []Model.Order{}
		seller.Role = "seller"

		//Validation for all field like email, password, mobilenumber etc..
		u := make(map[string]interface{})
		u["dateOfBirth"] = seller.DateOfBirth
		u["email"] = seller.Email
		u["mobileNumber"] = seller.MobileNumber
		u["password"] = seller.Password
		if !validationForRegister(u, w) {
			return
		}

		// encrypting the updated password
		encryptedValue := helper.Encrypt([]byte(seller.Password), "Secret Key")
		str_encryptedVal := hex.EncodeToString(encryptedValue)
		seller.Password = str_encryptedVal

		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "firstName", Value: seller.FirstName},
					{Key: "lastName", Value: seller.LastName},
					{Key: "role", Value: seller.Role},
					{Key: "email", Value: seller.Email},
					{Key: "mobileNumber", Value: seller.MobileNumber},
					{Key: "password", Value: seller.Password},
					{Key: "dateOfBirth", Value: seller.DateOfBirth},
					{Key: "isPrime", Value: seller.IsPrime},
					{Key: "productsListedIds", Value: seller.ProductsListedIds},
					{Key: "productsSoldIds", Value: seller.ProductsSoldIds},
				},
			},
		}
		userUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "firstName", Value: seller.FirstName},
					{Key: "lastName", Value: seller.LastName},
					{Key: "role", Value: seller.Role},
					{Key: "email", Value: seller.Email},
					{Key: "mobileNumber", Value: seller.MobileNumber},
					{Key: "password", Value: seller.Password},
					{Key: "dateOfBirth", Value: seller.DateOfBirth},
					{Key: "isPrime", Value: seller.IsPrime},
				},
			},
		}

		collection := conn.ConnectDB("sellers")
		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&seller)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		userCollection := conn.ConnectDB("users")
		errr := userCollection.FindOneAndUpdate(context.TODO(), filter, userUpdate).Decode(&seller)
		if errr != nil {
			conn.GetError(err, w)
			return
		}

		// seller.Id = string(id)

		json.NewEncoder(w).Encode(seller)
	}
}

func DeleteSeller(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["id"]
		id,errrr := primitive.ObjectIDFromHex(ids)
		if errrr!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": id}
		collection := conn.ConnectDB("sellers")
		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		userCollection := conn.ConnectDB("users")
		userDeleteResult, err := userCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(deleteResult)
		json.NewEncoder(w).Encode(userDeleteResult)
	}

}

// Validation functions
func validationForRegister(seller map[string]interface{}, w http.ResponseWriter) bool {
	//validation for email
	email, ok := seller["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForEmail(email, w) {
		return false
	}

	//validation for password
	password, ok := seller["password"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForPassword(password, w) {
		return false
	}

	//validation for mobile number
	mobileNumber, ok := seller["mobileNumber"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForMobileNumber(mobileNumber, w) {
		return false
	}

	//validation of dob
	dob, ok := seller["dateOfBirth"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForDOB(dob, w) {
		return false
	}

	return true

}

func validationForEmail(email string, w http.ResponseWriter) bool {
	if len(email) == 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("Email is Empty!!")
		return false
	}
	if len(email) < 5 {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Invalid Email!!")
		json.NewEncoder(w).Encode("Invalid Email!!")
		return false
	}
	if string(email[len(email)-4:]) != ".com" {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Invalid Email!!")
		json.NewEncoder(w).Encode("Invalid Email!!")
		return false
	}
	flag := 0
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			flag++
			break
		}
	}
	if flag == 0 {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Invalid Email!!")
		json.NewEncoder(w).Encode("Invalid Email!!")
		return false
	}
	return true
}

func validationForPassword(password string, w http.ResponseWriter) bool {
	if len(password) == 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("Password is Empty!!")
		return false
	}
	if len(password) < 8 {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Password should be more than 7 characters!!")
		json.NewEncoder(w).Encode("Password should be more than 7 characters!!")
		return false
	}
	cnt := 0
	for i := 0; i < len(password); i++ {
		c := password[i]
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			cnt++
		}
	}
	if cnt == 0 {
		w.WriteHeader(http.StatusConflict)
		fmt.Println("Password should contain atleast one super characters!!")
		json.NewEncoder(w).Encode("Password should contain atleast one super characters!!")
		return false
	}
	return true
}

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

func validationForDOB(dob string, w http.ResponseWriter) bool {
	if len(dob) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("DOB is Empty!!")
		return false
	}
	layout := "02/01/2006"
	date, err := time.Parse(layout, dob)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error parsing date: Invalid format!")
		return false
	}
	if date.After(time.Now()) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("Invalid DOB: Date is in the future!")
		return false
	}
	return true
}
