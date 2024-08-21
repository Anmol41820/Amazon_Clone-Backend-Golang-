package customer

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

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var customer []Model.Customer

	collection := conn.ConnectDB("customers")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Customer
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
		customer = append(customer, u)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(customer)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		customer := make(map[string]interface{})
		customer["role"] = "customer"
		customer["addresses"] = []Model.Address{}
		customer["productsPurchased"] = []Model.Order{}
		customer["isPrime"] = false

		err = json.Unmarshal([]byte(asString), &customer)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if customer == nil {
			http.Error(w, "Invalid customer data", http.StatusBadRequest)
			return
		}
		// delete(customer, "_id")

		//check wheather the register email or mobile number is already exist or not
		collection := conn.ConnectDB("customers")

		coll, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer coll.Close(context.TODO())

		for coll.Next(context.TODO()) {
			var existingcustomer Model.Customer
			err := coll.Decode(&existingcustomer)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if existingcustomer.Id == customer["_id"] {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Duplicate Id!!")
				return
			}
			if existingcustomer.Email == customer["email"].(string) || existingcustomer.MobileNumber == customer["mobileNumber"].(string) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Email or Mobile number already used!!")
				return
			}
		}
		if err := coll.Err(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		//Validation for all field like email, password, mobilenumber etc..
		if !validationForRegister(customer, w) {
			return
		}

		//register -> encry the password
		password, ok := customer["password"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encryptedValue := helper.Encrypt([]byte(password), "Secret Key")
		str_encryptedVal := hex.EncodeToString(encryptedValue)
		customer["password"] = str_encryptedVal
		customer["_id"] = primitive.NewObjectID()

		// collection := conn.ConnectDB("customers")
		result, err := collection.InsertOne(context.TODO(), customer)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		userCollection := conn.ConnectDB("users")
		userResult, err := userCollection.InsertOne(context.TODO(), customer)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//creating a new cart for the new customer
		newCart := make(map[string]interface{})
		newCart["_id"] = primitive.NewObjectID()
		newCart["customerId"] = customer["_id"]
		newCart["cartItems"] = []Model.CartItem{}
		newCart["numberOfProduct"] = 0
		newCart["totalAmount"] = 0
		cartCollection := conn.ConnectDB("carts")
		cartResult, err := cartCollection.InsertOne(context.TODO(), newCart)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//creating a new wishlist for the new customer
		newWishlist := make(map[string]interface{})
		newWishlist["_id"] = primitive.NewObjectID()
		newWishlist["customerId"] = customer["_id"]
		newWishlist["wishlistItems"] = []Model.WishlistItem{}
		newWishlist["numberOfProduct"] = 0
		wishlistCollection := conn.ConnectDB("wishlists")
		wishlistResult, err := wishlistCollection.InsertOne(context.TODO(), newWishlist)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//creating a new searchhistory for the new customer
		newSearchHistory := make(map[string]interface{})
		newSearchHistory["_id"] = primitive.NewObjectID()
		newSearchHistory["customerId"] = customer["_id"]
		newSearchHistory["searchText"] = []string{}
		searchHistoryCollection := conn.ConnectDB("searchHistories")
		searchHistoryResult, err := searchHistoryCollection.InsertOne(context.TODO(), newSearchHistory)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//creating a new productRecommendation for the new customer
		newProductRecommendation := make(map[string]interface{})
		newProductRecommendation["_id"] = primitive.NewObjectID()
		newProductRecommendation["customerId"] = customer["_id"]
		productRecommendationCollection := conn.ConnectDB("productRecommendations")
		productRecommendationResult, err := productRecommendationCollection.InsertOne(context.TODO(), newProductRecommendation)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//creating a new recentlyViewedProduct for the new customer
		newRecentlyViewedProduct := make(map[string]interface{})
		newRecentlyViewedProduct["_id"] = primitive.NewObjectID()
		newRecentlyViewedProduct["customerId"] = customer["_id"]
		recentlyViewedProductCollection := conn.ConnectDB("recentlyViewedProducts")
		recentlyViewedProductResult, err := recentlyViewedProductCollection.InsertOne(context.TODO(), newRecentlyViewedProduct)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(userResult)
		json.NewEncoder(w).Encode(cartResult)
		json.NewEncoder(w).Encode(wishlistResult)
		json.NewEncoder(w).Encode(searchHistoryResult)
		json.NewEncoder(w).Encode(productRecommendationResult)
		json.NewEncoder(w).Encode(recentlyViewedProductResult)
	}
}

func GetSingleCustomer(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("customers")

	var customer Model.Customer
	var params = mux.Vars(r)

	ids := params["id"]
	id, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&customer)
	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
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
		var customer Model.Customer

		filter := bson.M{"_id": id}

		_ = json.NewDecoder(r.Body).Decode(&customer)
		// customer.IsPrime = false
		// customer.Addresses = []Model.Address{}
		// customer.ProductsPurchased = []Model.Order{}
		customer.Role = "customer"

		//Validation for all field like email, password, mobilenumber etc..
		u := make(map[string]interface{})
		u["dateOfBirth"] = customer.DateOfBirth
		u["email"] = customer.Email
		u["mobileNumber"] = customer.MobileNumber
		u["password"] = customer.Password
		if !validationForRegister(u, w) {
			return
		}

		// encrypting the updated password
		encryptedValue := helper.Encrypt([]byte(customer.Password), "Secret Key")
		str_encryptedVal := hex.EncodeToString(encryptedValue)
		customer.Password = str_encryptedVal

		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "firstName", Value: customer.FirstName},
					{Key: "lastName", Value: customer.LastName},
					{Key: "role", Value: customer.Role},
					{Key: "email", Value: customer.Email},
					{Key: "mobileNumber", Value: customer.MobileNumber},
					{Key: "password", Value: customer.Password},
					{Key: "dateOfBirth", Value: customer.DateOfBirth},
					{Key: "isPrime", Value: customer.IsPrime},
					{Key: "productsPurchased", Value: customer.ProductsPurchased},
					{Key: "addresses", Value: customer.Addresses},
				},
			},
		}
		userUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "firstName", Value: customer.FirstName},
					{Key: "lastName", Value: customer.LastName},
					{Key: "role", Value: customer.Role},
					{Key: "email", Value: customer.Email},
					{Key: "mobileNumber", Value: customer.MobileNumber},
					{Key: "password", Value: customer.Password},
					{Key: "dateOfBirth", Value: customer.DateOfBirth},
					{Key: "isPrime", Value: customer.IsPrime},
				},
			},
		}

		collection := conn.ConnectDB("customers")
		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&customer)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		userCollection := conn.ConnectDB("users")
		errr := userCollection.FindOneAndUpdate(context.TODO(), filter, userUpdate).Decode(&customer)
		if errr != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(customer)
	}
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)
		ids := params["id"]
		id, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": id}
		collection := conn.ConnectDB("customers")
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

func AddMoneyInWallet(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the money from the body
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		moneyAmount := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &moneyAmount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//updating the customer wallet
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
		if int(moneyAmount["money"].(float64)) <0{
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w,"Can't be add negative money in the wallet!!")
			return
		}
		customer.Wallet += int(moneyAmount["money"].(float64))
		customerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wallet", Value: customer.Wallet},
				},
			},
		}
		customerResult,errrrr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
		if errrrr != nil {
			conn.GetError(errrrr, w)
			return
		}
		json.NewEncoder(w).Encode(customerResult)
		json.NewEncoder(w).Encode(customer)
	}
}

func WithdrawMoneyFromWallet(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the money from the body
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		moneyAmount := make(map[string]interface{})
		moneyAmount["withdrawAll"] = false
		err = json.Unmarshal([]byte(asString), &moneyAmount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//updating the customer wallet
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

		if moneyAmount["withdrawAll"].(bool){
			customer.Wallet = 0
		}else{
			if customer.Wallet < int(moneyAmount["money"].(float64)){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w,"Not enough Money in Wallet!!")
				return
			}else{
				customer.Wallet -= int(moneyAmount["money"].(float64))
			}
		}
		customerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wallet", Value: customer.Wallet},
				},
			},
		}
		customerResult,errrrr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
		if errrrr != nil {
			conn.GetError(errrrr, w)
			return
		}
		json.NewEncoder(w).Encode(customerResult)
		json.NewEncoder(w).Encode(customer)
	}
}

func BuyPrimeMembership(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the subscription from the body
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		subscription := make(map[string]interface{})
		// moneyAmount["withdrawAll"] = false
		err = json.Unmarshal([]byte(asString), &subscription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//updating the customer wallet
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

		if !customer.IsPrime{
			if int(subscription["months"].(float64))==1{
				if customer.Wallet >=299{
					customer.Wallet -= 299
					customer.IsPrime = true
					customer.PrimeExpireDate = time.Now().AddDate(0,1,0)
				}else{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Not enough Money in Wallet!!")
					return
				}
			}else if int(subscription["months"].(float64))==3{
				if customer.Wallet >=599{
					customer.Wallet -= 599
					customer.IsPrime = true
					customer.PrimeExpireDate = time.Now().AddDate(0,3,0)
				}else{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Not enough Money in Wallet!!")
					return
				}
			}else if int(subscription["months"].(float64))==6{
				if customer.Wallet >=899{
					customer.Wallet -= 899
					customer.IsPrime = true
					customer.PrimeExpireDate = time.Now().AddDate(0,6,0)
				}else{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Not enough Money in Wallet!!")
					return
				}
			}else if int(subscription["months"].(float64))==12{
				if customer.Wallet >=1099{
					customer.Wallet -= 1099
					customer.IsPrime = true
					customer.PrimeExpireDate = time.Now().AddDate(1,0,0)
				}else{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Not enough Money in Wallet!!")
					return
				}
			}
		}else{
			fmt.Fprintln(w,"You are already a prime member of amazon!!")
			return
		}
		customerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wallet", Value: customer.Wallet},
					{Key: "isPrime", Value: customer.IsPrime},
					{Key: "primeExpireDate", Value: customer.PrimeExpireDate},
				},
			},
		}
		customerResult,errrrr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
		if errrrr != nil {
			conn.GetError(errrrr, w)
			return
		}
		json.NewEncoder(w).Encode(customerResult)
		json.NewEncoder(w).Encode(customer)
	}
}

// Validation functions
func validationForRegister(customer map[string]interface{}, w http.ResponseWriter) bool {
	//validation for email
	email, ok := customer["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForEmail(email, w) {
		return false
	}

	//validation for password
	password, ok := customer["password"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForPassword(password, w) {
		return false
	}

	//validation for mobile number
	mobileNumber, ok := customer["mobileNumber"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForMobileNumber(mobileNumber, w) {
		return false
	}

	//validation of dob
	dob, ok := customer["dateOfBirth"].(string)
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
