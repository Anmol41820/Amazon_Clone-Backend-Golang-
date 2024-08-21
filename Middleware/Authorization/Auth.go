package auth

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	secretStr = "Secret Key"
	secretKey = []byte("Secret Key")
)

func ProtectedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Generic.SetupResponse(&w, r)

		//taken the id from the url
		vars := mux.Vars(r)
		ids, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Missing ID in URL")
			json.NewEncoder(w).Encode("Missing ID in URL")
			return
		}
		id,errrr := primitive.ObjectIDFromHex(ids)
		if errrr!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		collection := conn.ConnectDB("users")

		var user Model.User
		filter := bson.M{"_id": id}
		err := collection.FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			fmt.Println("i am here!!")
			conn.GetError(err, w)
			return
		}

		//verifing the token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenString := cookie.Value

		errr := VerifyToken(tokenString)
		if errr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Invalid token")
			json.NewEncoder(w).Encode("Invalid token")
			return
		}

		//verifing the perticuler user using email matching
		claims, err := RetriveEmailFromPayload(tokenString)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Println("Invalid token")
			json.NewEncoder(w).Encode("Invalid token")
            return
        }
		if claims["email"] != user.Email {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Println("User ID does not match token claims")
			json.NewEncoder(w).Encode("Unauthorized User!!")
            return
        }
		fmt.Println("Welcome to the protected area")
		next.ServeHTTP(w, r)
	})
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Invalid token")
	}
	return nil
}

func RetriveEmailFromPayload(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    if err != nil {
        return nil, err
    }
    if !token.Valid {
        return nil, fmt.Errorf("Invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("Invalid token claims")
    }

    return claims, nil
}