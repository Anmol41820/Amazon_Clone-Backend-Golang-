package review

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	"time"

	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetReviews(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var reviews []Model.Review
	//finding review in database
	collection := conn.ConnectDB("reviews")
	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Review
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		reviews = append(reviews, u)
	}
	json.NewEncoder(w).Encode(reviews)
}
func GetReviewsByProductId(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//tkaing the ids from the url
	var params = mux.Vars(r)
	ids := params["productId"]
	productId, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reviews []Model.Review
	//finding review in database
	collection := conn.ConnectDB("reviews")
	reviewFilter := bson.M{"productId": productId}
	curr, err := collection.Find(context.TODO(), reviewFilter)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Review
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		reviews = append(reviews, u)
	}
	json.NewEncoder(w).Encode(reviews)
}
func GetReviewsOfProductByMostRecent(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//tkaing the ids from the url
	var params = mux.Vars(r)
	ids := params["productId"]
	productId, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reviews []Model.Review
	//finding review in database
	collection := conn.ConnectDB("reviews")
	reviewFilter := bson.M{"productId": productId}
	reviewSortFilter := options.Find()
	reviewSortFilter.SetSort(bson.D{{"reviewDate", -1}})
	curr, err := collection.Find(context.TODO(), reviewFilter, reviewSortFilter)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Review
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		reviews = append(reviews, u)
	}
	json.NewEncoder(w).Encode(reviews)
}
func GetReviewsOfProductByTopReviews(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//tkaing the ids from the url
	var params = mux.Vars(r)
	ids := params["productId"]
	productId, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reviews []Model.Review
	//finding review in database
	collection := conn.ConnectDB("reviews")
	reviewFilter := bson.M{"productId": productId}
	reviewSortFilter := options.Find()
	reviewSortFilter.SetSort(bson.D{{"rating", -1}})
	curr, err := collection.Find(context.TODO(), reviewFilter, reviewSortFilter)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Review
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		reviews = append(reviews, u)
	}
	json.NewEncoder(w).Encode(reviews)
}

func AddReview(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the productid and customer id from url
		var params = mux.Vars(r)
		ids := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := params["productId"]
		productId, errrr := primitive.ObjectIDFromHex(id)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//takign the data for review through customer from body of postman
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var review Model.Review
		review.Id = primitive.NewObjectID()
		review.CustomerId = customerId
		review.ProductId = productId
		review.ReviewDate = time.Now()
		err = json.Unmarshal([]byte(asString), &review)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//fetching the product from product API to update the average rating and number of reviews
		var product Model.Product
		productCollection := conn.ConnectDB("products")
		productFilter := bson.M{"_id": productId}
		errr := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
		if errr != nil {
			conn.GetError(err, w)
			return
		}
		avergeRating := product.AverageRating
		newRating := float64(review.Rating)
		reviewCount := float64(product.NumberOfReviews)
		newAverageRating := (avergeRating*reviewCount + newRating)/(reviewCount+1.0)
		product.AverageRating = newAverageRating
		product.NumberOfReviews++
		productUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "averageRating", Value: product.AverageRating},
					{Key: "numberOfReviews", Value: product.NumberOfReviews},
				},
			},
		}
		productResult, errrr := productCollection.UpdateOne(context.TODO(), productFilter, productUpdate)
		if errrr != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(productResult)

		//connection to reviewss DB and inserting the review in it
		reviewCollection := conn.ConnectDB("reviews")
		reviewResult, err := reviewCollection.InsertOne(context.TODO(), review)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(reviewResult)
		json.NewEncoder(w).Encode(review)
	}
}
