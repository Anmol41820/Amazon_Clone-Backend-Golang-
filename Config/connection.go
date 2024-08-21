package connection

import (
    // Generic "Amazon_Server/Generic"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func ConnectDB(coll string) *mongo.Collection {
    // clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	clientOptions := options.Client().ApplyURI("mongodb+srv://Anmol41820:AI%4041820@cluster0.ezyeo6n.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB " + coll)
    collection := client.Database("Amazon_Clone").Collection(coll)

    return collection
}


func GetError(err error, w http.ResponseWriter) {

    log.Fatal(err.Error())
    var response = ErrorResponse{
        ErrorMessage: err.Error(),
        StatusCode:   http.StatusInternalServerError,
    }
    message, _ := json.Marshal(response)

    w.WriteHeader(response.StatusCode)
    w.Write(message)
}

var (
	secretStr = "Secret Key"
	secretKey = []byte("Secret Key")
)

func CreateToken(email string, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": id,
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}


type ErrorResponse struct {
    StatusCode   int    `json:"status"`
    ErrorMessage string `json:"message"`
}
