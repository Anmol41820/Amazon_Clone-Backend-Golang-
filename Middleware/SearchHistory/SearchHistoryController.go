package searchHistory

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"

	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSearchHistory(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var searchHistory []Model.SearchHistory
	//finding searchHistory in database
	collection := conn.ConnectDB("searchHistories")
	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.SearchHistory
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		searchHistory = append(searchHistory, u)
	}
	json.NewEncoder(w).Encode(searchHistory)
}

func GetSingleSearchHistoryByCustomerId(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	collection := conn.ConnectDB("searchHistories")
	var searchHistory Model.SearchHistory
	var params = mux.Vars(r)
	ids := params["id"]
	customerId, _ := primitive.ObjectIDFromHex(ids)
	filter := bson.M{"customerId": customerId}
	err := collection.FindOne(context.TODO(), filter).Decode(&searchHistory)
	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(searchHistory.SearchText)
}