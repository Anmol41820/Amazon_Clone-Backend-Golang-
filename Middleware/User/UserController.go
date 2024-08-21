package user

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

func GetUser(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var user []Model.User

	collection := conn.ConnectDB("users")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.User
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
		user = append(user, u)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := make(map[string]interface{})
		user["isPrime"] = false

		err = json.Unmarshal([]byte(asString), &user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user == nil {
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}
		// delete(user, "_id")

		//check wheather the register email or mobile number is already exist or not
		collection := conn.ConnectDB("users")

		coll, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer coll.Close(context.TODO())

		for coll.Next(context.TODO()) {
			var existingUser Model.User
			err := coll.Decode(&existingUser)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if existingUser.Id == user["_id"] {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Duplicate Id!!")
				return
			}
			if existingUser.Email == user["email"].(string) || existingUser.MobileNumber == user["mobileNumber"].(string) {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode("Email or Mobile number already used!!")
				return
			}
		}

		//Validation for all field like email, password, mobilenumber etc..
		if !validationForRegister(user, w) {
			return
		}

		//register -> encry the password
		password, ok := user["password"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		encryptedValue := helper.Encrypt([]byte(password), "Secret Key")
		str_encryptedVal := hex.EncodeToString(encryptedValue)
		user["password"] = str_encryptedVal

		user["_id"] = primitive.NewObjectID()
		// user["_id"] = tempId
		result, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		if user["role"] == "customer" {
			customerCollection := conn.ConnectDB("customers")
			// user["_id"] = tempId
			customerResult, err := customerCollection.InsertOne(context.TODO(), user)
			if err != nil {
				conn.GetError(err, w)
				return
			}
			//creating a new cart for the new customer
			newCart := make(map[string]interface{})
			newCart["_id"] = primitive.NewObjectID()
			newCart["customerId"] = user["_id"]
			newCart["CartItems"] = []Model.CartItem{}
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
			newWishlist["customerId"] = user["_id"]
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
			newSearchHistory["customerId"] = user["_id"]
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
			newProductRecommendation["customerId"] = user["_id"]
			productRecommendationCollection := conn.ConnectDB("productRecommendations")
			productRecommendationResult, err := productRecommendationCollection.InsertOne(context.TODO(), newProductRecommendation)
			if err != nil {
				conn.GetError(err, w)
				return
			}

			//creating a new recentlyViewedProduct for the new customer
			newRecentlyViewedProduct := make(map[string]interface{})
			newRecentlyViewedProduct["_id"] = primitive.NewObjectID()
			newRecentlyViewedProduct["customerId"] = user["_id"]
			recentlyViewedProductCollection := conn.ConnectDB("recentlyViewedProducts")
			recentlyViewedProductResult, err := recentlyViewedProductCollection.InsertOne(context.TODO(), newRecentlyViewedProduct)
			if err != nil {
				conn.GetError(err, w)
				return
			}

			json.NewEncoder(w).Encode(customerResult)
			json.NewEncoder(w).Encode(cartResult)
			json.NewEncoder(w).Encode(wishlistResult)
			json.NewEncoder(w).Encode(searchHistoryResult)
			json.NewEncoder(w).Encode(productRecommendationResult)
			json.NewEncoder(w).Encode(recentlyViewedProductResult)
		}
		
		if user["role"] == "seller" {
			sellerCollection := conn.ConnectDB("sellers")
			// user["_id"] = tempId
			sellerResult, err := sellerCollection.InsertOne(context.TODO(), user)
			if err != nil {
				conn.GetError(err, w)
				return
			}

			//creating a new report for the new seller
			newReport := make(map[string]interface{})
			newReport["_id"] = primitive.NewObjectID()
			newReport["sellerId"] = user["_id"]
			newReport["soldItems"] = []Model.SoldItem{}
			reportCollection := conn.ConnectDB("reports")
			reportResult, err := reportCollection.InsertOne(context.TODO(), newReport)
			if err != nil {
				conn.GetError(err, w)
				return
			}
			json.NewEncoder(w).Encode(sellerResult)
			json.NewEncoder(w).Encode(reportResult)
		}

		json.NewEncoder(w).Encode(result)
	}
}

func GetSingleUser(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("users")

	var user Model.User
	var params = mux.Vars(r)

	ids := params["id"]
	id, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
		var user Model.User

		filter := bson.M{"_id": id}

		_ = json.NewDecoder(r.Body).Decode(&user)

		//Validation for all field like email, password, mobilenumber etc..
		u := make(map[string]interface{})
		u["dateOfBirth"] = user.DateOfBirth
		u["email"] = user.Email
		u["mobileNumber"] = user.MobileNumber
		u["password"] = user.Password
		if !validationForRegister(u, w) {
			return
		}

		// encrypting the updated password
		encryptedValue := helper.Encrypt([]byte(user.Password), "Secret Key")
		str_encryptedVal := hex.EncodeToString(encryptedValue)
		user.Password = str_encryptedVal

		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "firstName", Value: user.FirstName},
					{Key: "lastName", Value: user.LastName},
					{Key: "role", Value: user.Role},
					{Key: "email", Value: user.Email},
					{Key: "mobileNumber", Value: user.MobileNumber},
					{Key: "password", Value: user.Password},
					{Key: "dateOfBirth", Value: user.DateOfBirth},
					{Key: "isPrime", Value: user.IsPrime},
				},
			},
		}

		collection := conn.ConnectDB("users")
		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		if user.Role == "customer" {
			customerCollection := conn.ConnectDB("customers")
			err := customerCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
			if err != nil {
				conn.GetError(err, w)
				return
			}
		}
		if user.Role == "seller" {
			sellerCollection := conn.ConnectDB("sellers")
			err := sellerCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
			if err != nil {
				conn.GetError(err, w)
				return
			}
		}

		// user.Id = string(id)

		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
		collection := conn.ConnectDB("users")
		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		customerCollection := conn.ConnectDB("customers")
		customerDeleteResult, err := customerCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		sellerCollection := conn.ConnectDB("sellers")
		sellerDeleteResult, err := sellerCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(deleteResult)
		json.NewEncoder(w).Encode(customerDeleteResult)
		json.NewEncoder(w).Encode(sellerDeleteResult)
	}

}

// Validation functions
func validationForRegister(user map[string]interface{}, w http.ResponseWriter) bool {
	//validation for email
	email, ok := user["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForEmail(email, w) {
		return false
	}

	//validation for password
	password, ok := user["password"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForPassword(password, w) {
		return false
	}

	//validation for mobile number
	mobileNumber, ok := user["mobileNumber"].(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !validationForMobileNumber(mobileNumber, w) {
		return false
	}

	//validation of dob
	dob, ok := user["dateOfBirth"].(string)
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
