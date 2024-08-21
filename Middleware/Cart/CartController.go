package cart

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

func GetCart(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var carts []Model.Cart

	collection := conn.ConnectDB("carts")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var cart Model.Cart
		err := curr.Decode(&cart)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}

		//taking the product using productId an updating the productInStock in wishlist
		productCollection := conn.ConnectDB("products")

		for i := 0; i < len(cart.CartItems); i++ {
			productId := cart.CartItems[i].ProductId
			var product Model.Product
			filterr := bson.M{"_id": productId}
			errr := productCollection.FindOne(context.TODO(), filterr).Decode(&product)
			if errr != nil {
				conn.GetError(err, w)
				return
			}

			if product.Quantity >= cart.CartItems[i].QuantityInCart {
				cart.CartItems[i].ProductInStock = true
				//to notify that product is in the stock....continue
			} else {
				cart.CartItems[i].ProductInStock = false
				if cart.CartItems[i].SelectedForBuying{
					cart.CartItems[i].SelectedForBuying = false
					cart.NumberOfProduct -= cart.CartItems[i].QuantityInCart
					cart.TotalAmount -= cart.CartItems[i].ProductPrice*cart.CartItems[i].QuantityInCart
				}
			}
		}
		filter := bson.M{"customerId": cart.CustomerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "cartItems", Value: cart.CartItems},
					{Key: "numberOfProduct", Value: cart.NumberOfProduct},
					{Key: "totalAmount", Value: cart.TotalAmount},
				},
			},
		}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		carts = append(carts, cart)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(carts)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
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
		myCart := GetMyCartByCustomerId(w, customerId)

		//takeing the productId, customerId and quantity of product from body
		productDetail := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &productDetail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if productDetail == nil {
			http.Error(w, "Invalid cart data", http.StatusBadRequest)
			return
		}
		// delete(cart, "_id")

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

		//loop through cartItems of myCart to check the current product exit or not
		flag := 1
		for i := 0; i < len(myCart.CartItems); i++ {
			if myCart.CartItems[i].ProductId == productIdFromBody {
				if myCart.CartItems[i].QuantityInCart + int(productDetail["quantity"].(float64)) > product.Quantity{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Out of Stock or Less Quantity Available!!")
					return 
				}
				myCart.CartItems[i].QuantityInCart += int(productDetail["quantity"].(float64))
				if myCart.CartItems[i].SelectedForBuying{
					myCart.NumberOfProduct += int(productDetail["quantity"].(float64))
					myCart.TotalAmount += int(productDetail["quantity"].(float64)) * product.SellingPrice
				}else{
					myCart.CartItems[i].SelectedForBuying = true
					myCart.NumberOfProduct += myCart.CartItems[i].QuantityInCart
					myCart.TotalAmount += myCart.CartItems[i].QuantityInCart*myCart.CartItems[i].ProductPrice
				}
				flag = 0
				break
			}
		}
		if flag == 1 {
			var newCartItem Model.CartItem
			newCartItem.ProductId, _ = primitive.ObjectIDFromHex(productDetail["productId"].(string))
			if int(productDetail["quantity"].(float64)) > product.Quantity{
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w,"Out of Stock or Less Quantity Available in Inventory!!")
				return 
			}
			newCartItem.ProductName = product.ProductName
			newCartItem.ProductPrice = product.SellingPrice
			newCartItem.ProductImage = product.ProductImages[0]
			newCartItem.QuantityInCart = int(productDetail["quantity"].(float64))
			newCartItem.SelectedForBuying = true
			newCartItem.ProductInStock = true
			myCart.CartItems = append(myCart.CartItems, newCartItem)
			myCart.TotalAmount += int(productDetail["quantity"].(float64)) * product.SellingPrice
			myCart.NumberOfProduct += int(productDetail["quantity"].(float64))
		}

		collection := conn.ConnectDB("carts")
		filter := bson.M{"customerId": customerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "cartItems", Value: myCart.CartItems},
					{Key: "numberOfProduct", Value: myCart.NumberOfProduct},
					{Key: "totalAmount", Value: myCart.TotalAmount},
				},
			},
		}
		_,err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(myCart)
	}
}

func GetCartByCustomerId(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("carts")

	var cart Model.Cart
	var params = mux.Vars(r)

	ids := params["id"]
	id, _ := primitive.ObjectIDFromHex(ids)
	filter := bson.M{"customerId": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&cart)
	if err != nil {
		conn.GetError(err, w)
		return
	}

	//taking the product using productId an updating the productInStock in wishlist
	productCollection := conn.ConnectDB("products")

	for i := 0; i < len(cart.CartItems); i++ {
		productId := cart.CartItems[i].ProductId
		var product Model.Product
		filterr := bson.M{"_id": productId}
		errr := productCollection.FindOne(context.TODO(), filterr).Decode(&product)
		if errr != nil {
			conn.GetError(err, w)
			return
		}

		if product.Quantity >= cart.CartItems[i].QuantityInCart {
			cart.CartItems[i].ProductInStock = true
			//to notify that product is in the stock....continue
		} else {
			cart.CartItems[i].ProductInStock = false
			if cart.CartItems[i].SelectedForBuying{
				cart.CartItems[i].SelectedForBuying = false
				cart.NumberOfProduct -= cart.CartItems[i].QuantityInCart
				cart.TotalAmount -= cart.CartItems[i].ProductPrice*cart.CartItems[i].QuantityInCart
			}
		}
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "cartItems", Value: cart.CartItems},
				{Key: "numberOfProduct", Value: cart.NumberOfProduct},
				{Key: "totalAmount", Value: cart.TotalAmount},
			},
		},
	}
	_, errr := collection.UpdateOne(context.TODO(), filter, update)
	if errr != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func RemoveAllItemsFromCart(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		collection := conn.ConnectDB("carts")

		var cart Model.Cart
		var params = mux.Vars(r)

		ids := params["id"]
		id, _ := primitive.ObjectIDFromHex(ids)
		filter := bson.M{"customerId": id}
		err := collection.FindOne(context.TODO(), filter).Decode(&cart)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		if len(cart.CartItems)==0{
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w,"Cart is already empty!!")
			return
		}
		var emptyCartItems []Model.CartItem
		tempZero := 0
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "cartItems", Value: emptyCartItems},
					{Key: "numberOfProduct", Value: tempZero},
					{Key: "totalAmount", Value: tempZero},
				},
			},
		}
		cartResult, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(cart)
		json.NewEncoder(w).Encode(cartResult)
	}

}

func GetMyCartByCustomerId(w http.ResponseWriter, id primitive.ObjectID) Model.Cart {
	collection := conn.ConnectDB("carts")
	var cart Model.Cart
	filter := bson.M{"customerId": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&cart)

	if err != nil {
		conn.GetError(err, w)
	}
	return cart
}

func ToggleBuyingProduct(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Invalid cart data", http.StatusBadRequest)
			return
		}
		// delete(cart, "_id")

		customerId, errrr := primitive.ObjectIDFromHex(productDetail["customerId"].(string))
		if errrr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		myCart := GetMyCartByCustomerId(w, customerId)
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

		//loop through cartItems of myCart to check the current product exit or not
		for i := 0; i < len(myCart.CartItems); i++ {
			if myCart.CartItems[i].ProductId == productIdFromBody {
				if myCart.CartItems[i].SelectedForBuying{
					myCart.TotalAmount -= product.SellingPrice * myCart.CartItems[i].QuantityInCart 
					myCart.NumberOfProduct -= myCart.CartItems[i].QuantityInCart
					myCart.CartItems[i].SelectedForBuying = !myCart.CartItems[i].SelectedForBuying 
				}else{
					if myCart.CartItems[i].ProductInStock{
						myCart.TotalAmount += product.SellingPrice * myCart.CartItems[i].QuantityInCart
						myCart.NumberOfProduct += myCart.CartItems[i].QuantityInCart
						myCart.CartItems[i].SelectedForBuying = !myCart.CartItems[i].SelectedForBuying
					}else{
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w,"Product Out of Stock, Unable to toggle to Buy!!")
						return
					}
				}
				break
			}
		}

		collection := conn.ConnectDB("carts")
		filter := bson.M{"customerId": customerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "cartItems", Value: myCart.CartItems},
					{Key: "numberOfProduct", Value: myCart.NumberOfProduct},
					{Key: "totalAmount", Value: myCart.TotalAmount},
				},
			},
		}
		_,err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(myCart)
	}
}

func IncreasingProductQuantity(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Invalid cart data", http.StatusBadRequest)
			return
		}
		// delete(cart, "_id")

		customerId, errrr := primitive.ObjectIDFromHex(productDetail["customerId"].(string))
		if errrr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		myCart := GetMyCartByCustomerId(w, customerId)
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

		//loop through cartItems of myCart to check the current product exit or not
		for i := 0; i < len(myCart.CartItems); i++ {
			if myCart.CartItems[i].ProductId == productIdFromBody {
				if product.Quantity >= myCart.CartItems[i].QuantityInCart + 1{
					myCart.CartItems[i].QuantityInCart ++
					if myCart.CartItems[i].SelectedForBuying{
						myCart.NumberOfProduct += 1
						myCart.TotalAmount += product.SellingPrice
					}
					break
				}else{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Out of Stock or Less Quantity Available in Inventory!!")
					return
				}
			}
		}

		collection := conn.ConnectDB("carts")
		filter := bson.M{"customerId": customerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "cartItems", Value: myCart.CartItems},
					{Key: "numberOfProduct", Value: myCart.NumberOfProduct},
					{Key: "totalAmount", Value: myCart.TotalAmount},
				},
			},
		}
		_,err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(myCart)
	}
}


func DecreasingProductQuantity(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Invalid cart data", http.StatusBadRequest)
			return
		}
		// delete(cart, "_id")

		customerId, errrr := primitive.ObjectIDFromHex(productDetail["customerId"].(string))
		if errrr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		myCart := GetMyCartByCustomerId(w, customerId)
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

		//loop through cartItems of myCart to check the current product exit or not
		for i := 0; i < len(myCart.CartItems); i++ {
			if myCart.CartItems[i].ProductId == productIdFromBody {
				if myCart.CartItems[i].QuantityInCart -1 == 0{
					if myCart.CartItems[i].SelectedForBuying{
						myCart.NumberOfProduct --
						myCart.TotalAmount -= product.SellingPrice
					}
					//deleting the cartItem from myCart whose quantity become 0
					myCart.CartItems = append(myCart.CartItems[:i], myCart.CartItems[i+1:]...)
					break
				}
				myCart.CartItems[i].QuantityInCart --
				if myCart.CartItems[i].SelectedForBuying{
					myCart.NumberOfProduct --
					myCart.TotalAmount -= product.SellingPrice
				}
				break
			}
		}

		collection := conn.ConnectDB("carts")
		filter := bson.M{"customerId": customerId}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "cartItems", Value: myCart.CartItems},
					{Key: "numberOfProduct", Value: myCart.NumberOfProduct},
					{Key: "totalAmount", Value: myCart.TotalAmount},
				},
			},
		}
		_,err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		json.NewEncoder(w).Encode(myCart)
	}
}