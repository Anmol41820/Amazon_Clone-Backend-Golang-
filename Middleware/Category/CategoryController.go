package category

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


func GetCategory(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var category []Model.Category

	collection := conn.ConnectDB("categories")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Category
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
		category = append(category, u)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(category)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// var params = mux.Vars(r)
		// ids := params["id"]
		// id,errrr := primitive.ObjectIDFromHex(ids)
		// if errrr!=nil{
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }
		category := make(map[string]interface{})
		// category["averageRating"] = 0.0

		err = json.Unmarshal([]byte(asString), &category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if category == nil {
			http.Error(w, "Invalid category data", http.StatusBadRequest)
			return
		}
		// delete(category, "_id")

		collection := conn.ConnectDB("categories")
		result, err := collection.InsertOne(context.TODO(), category)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(result)
	}
}

func GetSingleCategory(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("categories")

	var category Model.Category
	var params = mux.Vars(r)

	ids := params["categoryId"]
	id,_ := primitive.ObjectIDFromHex(ids)
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&category)
	
	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(category)
}


func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["categoryId"]
		id,_ := primitive.ObjectIDFromHex(ids)

		filter := bson.M{"_id": id}
		collection := conn.ConnectDB("categories")
		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(deleteResult)
	}

}