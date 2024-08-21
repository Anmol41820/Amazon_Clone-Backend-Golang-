package forgotPassword

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	helper "Amazon_Server/Helper"

	"Amazon_Server/Model"
	"context"
	"fmt"

	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"net/http"
	"go.mongodb.org/mongo-driver/bson"
)


func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)

		var forgotPassword map[string]interface{}
		json.Unmarshal([]byte(asString), &forgotPassword)

		email, ok := forgotPassword["email"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		//validation for email
		if !validationForEmail(email,w){
			return 
		}

		collection := conn.ConnectDB("users")
		filter := bson.M{"email": email}
		var user Model.User
		errr := collection.FindOne(context.TODO(), filter).Decode(&user)
		if errr != nil {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode("User not registered, Please Register first!!")
			return
		}

		if err != nil {
			conn.GetError(err, w)
			return
		}

		newPass, ok := forgotPassword["newPassword"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		//validation for new password
		if !validationForPassword(newPass,w){
			return 
		}
		//change the password and encrypted it
		user.Password = newPass
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
		errrr := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
		if errrr != nil {
			conn.GetError(err, w)
			return
		}


		if user.Role == "customer"{
			customerCollection := conn.ConnectDB("customers")
			var customer Model.Customer
			errr := customerCollection.FindOne(context.TODO(), filter).Decode(&customer)
			if errr != nil {
				w.WriteHeader(http.StatusBadGateway)
				json.NewEncoder(w).Encode("Customer not registered, Please Register first!!")
				return
			}

			customer.Password = newPass
			customerEncryptedValue := helper.Encrypt([]byte(customer.Password), "Secret Key")
			customer_str_encryptedVal := hex.EncodeToString(customerEncryptedValue)
			customer.Password = customer_str_encryptedVal

			update := bson.D{
				{
					Key: "$set", Value: bson.D{
						{Key: "password", Value: user.Password},
					},
				},
			}
			errrr := customerCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&customer)
			if errrr != nil {
				conn.GetError(err, w)
				return
			}
		}

		if user.Role == "seller"{
			sellerCollection := conn.ConnectDB("sellers")
			var seller Model.Seller
			errr := sellerCollection.FindOne(context.TODO(), filter).Decode(&seller)
			if errr != nil {
				w.WriteHeader(http.StatusBadGateway)
				json.NewEncoder(w).Encode("seller not registered, Please Register first!!")
				return
			}

			seller.Password = newPass
			sellerEncryptedValue := helper.Encrypt([]byte(seller.Password), "Secret Key")
			seller_str_encryptedVal := hex.EncodeToString(sellerEncryptedValue)
			seller.Password = seller_str_encryptedVal

			update := bson.D{
				{
					Key: "$set", Value: bson.D{
						{Key: "password", Value: seller.Password},
					},
				},
			}
			errrr := sellerCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&seller)
			if errrr != nil {
				conn.GetError(err, w)
				return
			}
		}

		fmt.Fprintln(w,"Password changed successfully!!")
		// json.NewEncoder(w).Encode(user)
	}
}


//Validation functions
func validationForEmail(email string, w http.ResponseWriter) bool{
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

func validationForPassword(password string,w http.ResponseWriter) bool{
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