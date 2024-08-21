package login

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	helper "Amazon_Server/Helper"

	"Amazon_Server/Model"
	"context"
	"fmt"
	"time"

	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"net/http"
	"go.mongodb.org/mongo-driver/bson"
)


func Login(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)

		var login map[string]interface{}
		json.Unmarshal([]byte(asString), &login)

		email, ok := login["email"].(string)
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

		hex_pass, _ := hex.DecodeString(user.Password)
		decryptedValue := helper.Decrypt(hex_pass, "Secret Key")

		if err != nil {
			conn.GetError(err, w)
			return
		}

		pass, ok := login["password"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if pass == string(decryptedValue[:]) {
			fmt.Println("Valid User!! for ", email)
			tokenString, err := conn.CreateToken(email,user.Id.Hex())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cookie := &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/", // Ensure the cookie is valid for the entire site
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode, // Adjust SameSite policy as needed
				// Secure:   true, // Uncomment if using HTTPS
			}
			http.SetCookie(w, cookie)
			json.NewEncoder(w).Encode("Welcome!!")
			fmt.Fprintln(w,user.FirstName + ", You have login as "+ user.Role+"!!")
		} else {
			fmt.Println("Invalid User!! for ", email)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid User or Wrong password!!")
		}

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