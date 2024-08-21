package wishlist

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	// wishlist "Amazon_Server/Middleware/Wishlist"

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

func GetWishlist(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var wishlists []Model.Wishlist

	collection := conn.ConnectDB("wishlists")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var wishlist Model.Wishlist
		err := curr.Decode(&wishlist)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}

		//taking the product using productId an updating the productInStock in wishlist
		productCollection := conn.ConnectDB("products")

		for i := 0; i < len(wishlist.WishlistItems); i++ {
			productId := wishlist.WishlistItems[i].ProductId
			var product Model.Product
			filterr := bson.M{"_id": productId}
			errr := productCollection.FindOne(context.TODO(), filterr).Decode(&product)
			if errr != nil {
				conn.GetError(err, w)
				return
			}

			if product.Quantity > 0 {
				wishlist.WishlistItems[i].ProductInStock = true
			} else {
				wishlist.WishlistItems[i].ProductInStock = false
			}
		}
		filter := bson.M{"customerId": wishlist.CustomerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wishlistItems", Value: wishlist.WishlistItems},
				},
			},
		}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		wishlists = append(wishlists, wishlist)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(wishlists)
}

func AddToWishlist(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var params = mux.Vars(r)
		ids := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		myWishlist := GetMyWishlistByCustomerId(w, customerId)

		//takeing the productId, customerId and quantity of product from body
		productDetail := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &productDetail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if productDetail == nil {
			http.Error(w, "Invalid wishlist data", http.StatusBadRequest)
			return
		}
		// delete(wishlist, "_id")

		//getting productId from body of postman
		productIdFromBody, _ := primitive.ObjectIDFromHex(productDetail["productId"].(string))

		//taking the product using productId an updating the quantity if condition satisfies
		productCollection := conn.ConnectDB("products")
		var product Model.Product
		filterr := bson.M{"_id": productIdFromBody}
		errr := productCollection.FindOne(context.TODO(), filterr).Decode(&product)
		if errr != nil {
			conn.GetError(err, w)
			return
		}
		// fmt.Println("length: ",myWishlist.WishlistItems[0].ProductId)

		//loop through wishlistItems of mywishlist to check the current product exit or not
		flag := 1
		for _, value := range myWishlist.WishlistItems {
			if value.ProductId == productIdFromBody {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintln(w, "Product is already in your Wishlist!!")
				return
			}
		}
		if flag == 1 {
			var newWishlistItem Model.WishlistItem
			newWishlistItem.ProductId, _ = primitive.ObjectIDFromHex(productDetail["productId"].(string))
			newWishlistItem.ProductName = product.ProductName
			newWishlistItem.ProductPrice = product.SellingPrice
			newWishlistItem.ProductImage = product.ProductImages[0]
			if product.Quantity > 0 {
				newWishlistItem.ProductInStock = true
			} else {
				newWishlistItem.ProductInStock = false
			}
			myWishlist.WishlistItems = append(myWishlist.WishlistItems, newWishlistItem)
		}
		myWishlist.NumberOfProduct = len(myWishlist.WishlistItems)

		collection := conn.ConnectDB("wishlists")
		filter := bson.M{"customerId": customerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wishlistItems", Value: myWishlist.WishlistItems},
					{Key: "numberOfProduct", Value: myWishlist.NumberOfProduct},
				},
			},
		}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(myWishlist)
	}
}

func GetWishlistByCustomerId(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("wishlists")

	var wishlist Model.Wishlist
	var params = mux.Vars(r)

	ids := params["id"]
	id, _ := primitive.ObjectIDFromHex(ids)
	filter := bson.M{"customerId": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&wishlist)
	if err != nil {
		conn.GetError(err, w)
		return
	}

	//taking the product using productId an updating the productInStock in wishlist
	productCollection := conn.ConnectDB("products")

	for i := 0; i < len(wishlist.WishlistItems); i++ {
		productId := wishlist.WishlistItems[i].ProductId
		var product Model.Product
		filterr := bson.M{"_id": productId}
		errr := productCollection.FindOne(context.TODO(), filterr).Decode(&product)
		if errr != nil {
			conn.GetError(err, w)
			return
		}

		if product.Quantity > 0 {
			wishlist.WishlistItems[i].ProductInStock = true
		} else {
			wishlist.WishlistItems[i].ProductInStock = false
		}
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "wishlistItems", Value: wishlist.WishlistItems},
			},
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(wishlist)
}

func RemoveAllItemsFromWishlist(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		collection := conn.ConnectDB("wishlists")

		var wishlist Model.Wishlist
		var params = mux.Vars(r)

		ids := params["id"]
		id, _ := primitive.ObjectIDFromHex(ids)
		filter := bson.M{"customerId": id}
		err := collection.FindOne(context.TODO(), filter).Decode(&wishlist)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		if len(wishlist.WishlistItems)==0{
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w,"Wishlist is already empty!!")
			return
		}
		var emptyWishlistItems []Model.WishlistItem
		tempZero := 0
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wishlistItems", Value: emptyWishlistItems},
					{Key: "numberOfProduct", Value: tempZero},
				},
			},
		}
		wishlistResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(wishlist)
		json.NewEncoder(w).Encode(wishlistResult)
	}

}

func GetMyWishlistByCustomerId(w http.ResponseWriter, id primitive.ObjectID) Model.Wishlist {
	collection := conn.ConnectDB("wishlists")
	var wishlist Model.Wishlist
	filter := bson.M{"customerId": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&wishlist)

	if err != nil {
		conn.GetError(err, w)
	}
	return wishlist
}

func RemoveProductFromWishlist(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//takeing the productId, customerId and quantity of product from body
		productDetail := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &productDetail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if productDetail == nil {
			http.Error(w, "Invalid wishlist data", http.StatusBadRequest)
			return
		}
		// delete(wishlist, "_id")

		customerId, errrr := primitive.ObjectIDFromHex(productDetail["customerId"].(string))
		if errrr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		myWishlist := GetMyWishlistByCustomerId(w, customerId)
		productIdFromBody, _ := primitive.ObjectIDFromHex(productDetail["productId"].(string))

		//taking the product using productId an updating the quantity if condition satisfies
		productCollection := conn.ConnectDB("products")
		var product Model.Product
		filterr := bson.M{"_id": productIdFromBody}
		errr := productCollection.FindOne(context.TODO(), filterr).Decode(&product)
		if errr != nil {
			conn.GetError(err, w)
			return
		}

		//loop through wishlistItems of mywishlist to check the current product exit or not
		for i := 0; i < len(myWishlist.WishlistItems); i++ {
			if myWishlist.WishlistItems[i].ProductId == productIdFromBody {
				myWishlist.NumberOfProduct--
				myWishlist.WishlistItems = append(myWishlist.WishlistItems[:i], myWishlist.WishlistItems[i+1:]...)
				break
			}
		}

		collection := conn.ConnectDB("wishlists")
		filter := bson.M{"customerId": customerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "wishlistItems", Value: myWishlist.WishlistItems},
					{Key: "numberOfProduct", Value: myWishlist.NumberOfProduct},
				},
			},
		}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(myWishlist)
	}
}
