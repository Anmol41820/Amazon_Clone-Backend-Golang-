package recentlyViewedProduct

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"

	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRecentlyViewedProductByCustomerId(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var recentlyViewedProduct Model.RecentlyViewedProduct

	//fetching the customerId from the url
	var params = mux.Vars(r)
	idx := params["id"]
	customerId, _ := primitive.ObjectIDFromHex(idx)

	//finding recentlyViewedProduct in database
	collection := conn.ConnectDB("recentlyViewedProducts")
	filter := bson.M{"customerId": customerId}
	errr := collection.FindOne(context.TODO(), filter).Decode(&recentlyViewedProduct)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}

	json.NewEncoder(w).Encode(recentlyViewedProduct)
}

