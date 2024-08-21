package order

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
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var order []Model.Order

	collection := conn.ConnectDB("orders")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Order
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		order = append(order, u)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(order)
}

func BuyNow(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var order Model.Order
		err = json.Unmarshal([]byte(asString), &order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order.Id = primitive.NewObjectID()

		//fetching hte customer from customer database
		var customer Model.Customer
		customerCollection := conn.ConnectDB("customers")
		customerFilter := bson.M{"_id": order.CustomerId}
		errr := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
		if errr != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		//finding the default address if address is not provided by customer
		for k:=0;k<len(customer.Addresses);k++{
			if customer.Addresses[k].IsDefault{
				order.ShippingAddress = customer.Addresses[k]
				break
			}
		}

		productCollection := conn.ConnectDB("products")
		for i := 0; i < len(order.ProductIds); i++ {

			//checking the out of stock product
			var product Model.Product
			productFilter := bson.M{"_id": order.ProductIds[i]}
			err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			if product.Quantity < order.OrderQuantitys[i] {
				w.WriteHeader(http.StatusConflict)
				fmt.Fprintln(w, order.ProductNames[i]+" is Out of Stock, Please Reduce the quantity of the product!!")
				return
			}
			//calculating the total price
			order.PriceDetail.ListPrice += product.MaxRetailPrice * order.OrderQuantitys[i]
			order.PriceDetail.SellingPrice += product.SellingPrice * order.OrderQuantitys[i]
		}
		if order.PriceDetail.SellingPrice < 200 && !customer.IsPrime {
			order.PriceDetail.DeliveryCharge = 99
			order.PriceDetail.TotalAmount = order.PriceDetail.DeliveryCharge + order.PriceDetail.SellingPrice
		} else {
			order.PriceDetail.TotalAmount = order.PriceDetail.SellingPrice
		}
		order.TotalAmount = order.PriceDetail.TotalAmount

		//inserting the order in order API
		orderCollection := conn.ConnectDB("orders")
		result, err := orderCollection.InsertOne(context.TODO(), order)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(order.Id)
		json.NewEncoder(w).Encode(order.PriceDetail)
	}
}

func ContinueTOPayment(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the order from orderId get from url
		var order Model.Order
		var params = mux.Vars(r)
		ids := params["orderId"]
		orderId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&order)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}

		//taking the paymentDetial from the body of the postman
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var paymentDetails Model.Payment
		err = json.Unmarshal([]byte(asString), &paymentDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order.OrderedDate = time.Now()

		//payment
		order.PaymentDetails.Id = primitive.NewObjectID()
		order.PaymentDetails.PaymentMethod = paymentDetails.PaymentMethod
		order.PaymentDetails.TotalAmount = order.TotalAmount
		order.PaymentDetails.UserId = paymentDetails.UserId
		order.PaymentDetails.CardNumber = paymentDetails.CardNumber
		order.PaymentDetails.CardExpiryDate = paymentDetails.CardExpiryDate
		order.PaymentDetails.NameOnCard = paymentDetails.NameOnCard
		order.PaymentDetails.UpiId = paymentDetails.UpiId

		//fatching the customer
		var customer Model.Customer
		customerCollection := conn.ConnectDB("customers")
		customerFilter := bson.M{"_id": order.CustomerId}
		errrrr := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
		if errrrr != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		//payment method check
		if order.PaymentDetails.PaymentMethod == "Wallet" {
			if customer.Wallet >= order.TotalAmount {
				customer.Wallet -= order.TotalAmount
				customerUpdate := bson.D{
					{
						Key: "$set", Value: bson.D{
							{Key: "wallet", Value: customer.Wallet},
						},
					},
				}
				_, errr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
				if errr != nil {
					conn.GetError(err, w)
					return
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Not Enough Amount in your Wallet!!")
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Only Wallet Payment Available!!")
			return
		}

		//reducing the quantity from the inventry
		for i := 0; i < len(order.ProductIds); i++ {
			order.DeliveredDates = append(order.DeliveredDates, []time.Time{})
			order.DeliveredDates[i] = append(order.DeliveredDates[i], time.Now())
			if customer.IsPrime {
				order.DeliveredDates[i] = append(order.DeliveredDates[i], time.Now().AddDate(0, 0, 1))
			} else {
				order.DeliveredDates[i] = append(order.DeliveredDates[i], time.Now().AddDate(0, 0, 3))
			}
			order.Status = append(order.Status, []string{})
			order.Status[i] = append(order.Status[i], "Order Confirmed")

			var product Model.Product
			productCollection := conn.ConnectDB("products")
			productFilter := bson.M{"_id": order.ProductIds[i]}
			err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			product.Quantity -= order.OrderQuantitys[i]
			product.UnitsSold += order.OrderQuantitys[i]
			productUpdate := bson.D{
				{
					Key: "$set", Value: bson.D{
						{Key: "quantity", Value: product.Quantity},
						{Key: "unitsSold", Value: product.UnitsSold},
					},
				},
			}
			_, errr := productCollection.UpdateOne(context.TODO(), productFilter, productUpdate)
			if errr != nil {
				conn.GetError(err, w)
				return
			}

			//inserting the deliver order in deliverOrder model
			var deliverOrder Model.DeliverOrder
			deliverOrder.Id = primitive.NewObjectID()
			deliverOrder.CustomerId = order.CustomerId
			deliverOrder.ProductId = order.ProductIds[i]
			deliverOrder.OrderId = order.Id
			deliverOrder.CustomerName = order.ShippingAddress.FullName
			deliverOrder.CustomerAddress = order.ShippingAddress
			deliverOrder.ExpectedDeliveryDate = order.DeliveredDates[i][1]

			deliverOrderCollection := conn.ConnectDB("deliverOrders")
			_, errrr := deliverOrderCollection.InsertOne(context.TODO(), deliverOrder)
			if errrr != nil {
				conn.GetError(err, w)
				return
			}
		}

		orderUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "deliveredDates", Value: order.DeliveredDates},
					{Key: "orderedDate", Value: order.OrderedDate},
					{Key: "paymentDetails", Value: order.PaymentDetails},
					{Key: "priceDetail", Value: order.PriceDetail},
					{Key: "status", Value: order.Status},
					{Key: "totalAmount", Value: order.TotalAmount},
					{Key: "shippingAddress", Value: order.ShippingAddress},
					{Key: "customerId", Value: order.CustomerId},
					{Key: "orderQuantitys", Value: order.OrderQuantitys},
					{Key: "productIds", Value: order.ProductIds},
					{Key: "productNames", Value: order.ProductNames},
				},
			},
		}
		result, err := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
		json.NewEncoder(w).Encode(order)
	}
}

func GetOrdersByCustomerId(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	collection := conn.ConnectDB("orders")

	var order Model.Order
	var params = mux.Vars(r)

	ids := params["id"]
	id, errrr := primitive.ObjectIDFromHex(ids)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter := bson.M{"customerId": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func ReplaceOrder(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	// if(!conn.ProtectedHandler(w,r)){return}

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")

		//taking the issue from customer from body of postman
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		issueDetails := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &issueDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//gettign the orderid and productid from the url
		var params = mux.Vars(r)
		ids := params["orderId"]
		orderId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idx := params["productId"]
		productId, errrr := primitive.ObjectIDFromHex(idx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var oldOrder Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&oldOrder)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}

		for i := 0; i < len(oldOrder.ProductIds); i++ {
			if oldOrder.ProductIds[i] == productId {
				//taking the product which we need to replace
				var product Model.Product
				productCollection := conn.ConnectDB("products")
				productFilter := bson.M{"_id": productId}
				err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
				if err != nil {
					conn.GetError(err, w)
					return
				}

				if !product.ReplacePolicy {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "No Replace Policy for this product!!")
					return
				}

				if oldOrder.Status[i][len(oldOrder.Status[i])-1] != "Order Delivered" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "Please Replace when Order Delivered or Only One time Replace Policy!!")
					return
				} else {
					n := len(oldOrder.DeliveredDates[i])
					if time.Now().Before(oldOrder.DeliveredDates[i][n-1].AddDate(0, 0, 7)) {
						//checking the product is out of stock or not to replace
						if product.Quantity < oldOrder.OrderQuantitys[i] {
							w.WriteHeader(http.StatusBadRequest)
							fmt.Fprintln(w, "Product Unavailable, Please try for Return!!")
							return
						}
						//updating the product quantity in product model
						product.Quantity -= oldOrder.OrderQuantitys[i]
						productUpdate := bson.D{
							{
								Key: "$set", Value: bson.D{
									{Key: "quantity", Value: product.Quantity},
								},
							},
						}
						_, errrrr := productCollection.UpdateOne(context.TODO(), productFilter, productUpdate)
						if errrrr != nil {
							conn.GetError(errrrr, w)
							return
						}

						//write the code below there to replace
						oldOrder.Status[i] = append(oldOrder.Status[i], "Replace Confirmed")
						oldOrder.DeliveredDates[i] = append(oldOrder.DeliveredDates[i], time.Now())
						oldOrder.DeliveredDates[i] = append(oldOrder.DeliveredDates[i], time.Now().AddDate(0, 0, 3))

						//inserting the issue of the customer in replace model
						var replaceOrder Model.ReplaceOrder
						replaceOrder.Id = primitive.NewObjectID()
						replaceOrder.CustomerId = oldOrder.CustomerId
						replaceOrder.ProductId = oldOrder.ProductIds[i]
						replaceOrder.OrderId = oldOrder.Id
						replaceOrder.CustomerName = oldOrder.ShippingAddress.FullName
						replaceOrder.CustomerAddress = oldOrder.ShippingAddress
						nn := len(oldOrder.DeliveredDates[i])
						replaceOrder.ExpectedReplaceDate = oldOrder.DeliveredDates[i][nn-1]

						if issueDetails["isDamage"] == true {
							replaceOrder.IsDamage = true
						} else if issueDetails["dontLikeDueToColorOrSize"] == true {
							replaceOrder.DontLikeDueToColorOrSize = true
						}
						replaceOrderCollection := conn.ConnectDB("replaceOrders")
						_, err := replaceOrderCollection.InsertOne(context.TODO(), replaceOrder)
						if err != nil {
							conn.GetError(err, w)
							return
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, "Replace Expires!!")
						return
					}
				}
			}
		}

		orderUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "status", Value: oldOrder.Status},
					{Key: "deliveredDates", Value: oldOrder.DeliveredDates},
				},
			},
		}
		_, errrrr := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
		if errrrr != nil {
			conn.GetError(errrrr, w)
			return
		}

		json.NewEncoder(w).Encode(oldOrder)
	}
}

func ReturnOrder(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")
		//taking the issue from customer from body of postman
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		issueDetails := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &issueDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//gettign the orderid and productid from the url
		var params = mux.Vars(r)
		ids := params["orderId"]
		orderId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idx := params["productId"]
		productId, errrr := primitive.ObjectIDFromHex(idx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var oldOrder Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&oldOrder)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}

		for i := 0; i < len(oldOrder.ProductIds); i++ {
			if oldOrder.ProductIds[i] == productId {
				//taking the product which we need to replace
				var product Model.Product
				productCollection := conn.ConnectDB("products")
				productFilter := bson.M{"_id": productId}
				err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
				if err != nil {
					conn.GetError(err, w)
					return
				}

				if !product.ReturnPolicy {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "No Return Policy for this product!!")
					return
				}

				if oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Order Delivered" || oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Cancelled" || oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Replace Order Delivered"{
					n := len(oldOrder.DeliveredDates[i])
					if time.Now().Before(oldOrder.DeliveredDates[i][n-1].AddDate(0, 0, 7)) {
						//write the code below there to place
						oldOrder.Status[i] = append(oldOrder.Status[i], "Return Confirmed")
						oldOrder.DeliveredDates[i] = append(oldOrder.DeliveredDates[i], time.Now())
						oldOrder.DeliveredDates[i] = append(oldOrder.DeliveredDates[i], time.Now().AddDate(0, 0, 1))

						//inserting the issue of the customer in return model
						var returnOrder Model.ReturnOrder
						returnOrder.Id = primitive.NewObjectID()
						returnOrder.CustomerId = oldOrder.CustomerId
						returnOrder.ProductId = oldOrder.ProductIds[i]
						returnOrder.OrderId = oldOrder.Id
						returnOrder.CustomerName = oldOrder.ShippingAddress.FullName
						returnOrder.CustomerAddress = oldOrder.ShippingAddress
						nn := len(oldOrder.DeliveredDates[i])
						returnOrder.ExpectedReturnDate = oldOrder.DeliveredDates[i][nn-1]

						if issueDetails["isDamage"] == true {
							returnOrder.IsDamage = true
						} else if issueDetails["dontLikeDueToColorOrSize"] == true {
							returnOrder.DontLikeDueToColorOrSize = true
						}
						returnOrderCollection := conn.ConnectDB("returnOrders")
						_, err := returnOrderCollection.InsertOne(context.TODO(), returnOrder)
						if err != nil {
							conn.GetError(err, w)
							return
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, "Return Expires!!")
						return
					}
				} else {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "Please Refund when Order Delivered!!")
					return
				}
			}
		}

		orderUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "status", Value: oldOrder.Status},
					{Key: "deliveredDates", Value: oldOrder.DeliveredDates},
				},
			},
		}
		_, errrrr := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
		if errrrr != nil {
			conn.GetError(errrrr, w)
			return
		}

		json.NewEncoder(w).Encode(oldOrder)
	}
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["orderId"]
		orderId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idx := params["productId"]
		productId, errrr := primitive.ObjectIDFromHex(idx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idxx := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(idxx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var oldOrder Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		err := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&oldOrder)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		for i := 0; i < len(oldOrder.ProductIds); i++ {
			if oldOrder.ProductIds[i] == productId {
				//taking the product which we need to replace
				var product Model.Product
				productCollection := conn.ConnectDB("products")
				productFilter := bson.M{"_id": productId}
				err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
				if err != nil {
					conn.GetError(err, w)
					return
				}

				if oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Order Confirmed" {
					if time.Now().Before(oldOrder.OrderedDate.AddDate(0, 0, 1)) {
						//write the code below there to place
						oldOrder.Status[i] = append(oldOrder.Status[i], "Cancelled")
						n := len(oldOrder.DeliveredDates[i])
						oldOrder.DeliveredDates[i][n-1] = time.Now()

						//updating the product quantity in inventry
						product.Quantity += oldOrder.OrderQuantitys[i]
						product.UnitsSold -= oldOrder.OrderQuantitys[i]
						productUpdate := bson.D{
							{
								Key: "$set", Value: bson.D{
									{Key: "quantity", Value: product.Quantity},
									{Key: "unitsSold", Value: product.UnitsSold},
								},
							},
						}
						_, errr := productCollection.UpdateOne(context.TODO(), productFilter, productUpdate)
						if errr != nil {
							conn.GetError(errr, w)
							return
						}

						//deleting the deliverOrder request from deliverOrder model
						deliverOrderCollection := conn.ConnectDB("deliverOrders")
						deliverOrderFilter := bson.M{"productId": productId, "orderId": orderId}
						_, err := deliverOrderCollection.DeleteOne(context.TODO(), deliverOrderFilter)
						if err != nil {
							conn.GetError(err, w)
							return
						}

						//update the customer wallet
						if oldOrder.PaymentDetails.PaymentMethod == "Wallet" {
							var customer Model.Customer
							customerCollection := conn.ConnectDB("customers")
							customerFilter := bson.M{"_id": customerId}
							err := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
							if err != nil {
								conn.GetError(err, w)
								return
							}
							customer.Wallet += product.SellingPrice * oldOrder.OrderQuantitys[i]
							customerUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "wallet", Value: customer.Wallet},
									},
								},
							}
							_, errr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
							if errr != nil {
								conn.GetError(errr, w)
								return
							}
						}

					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, "Cancel Expires!!")
						return
					}
				} else if oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Replace Confirmed" {
					n := len(oldOrder.DeliveredDates[i])
					if time.Now().Before(oldOrder.DeliveredDates[i][n-1].AddDate(0, 0, -2)) {
						//write the code below there to place
						oldOrder.Status[i] = append(oldOrder.Status[i], "Cancelled")
						oldOrder.DeliveredDates[i][n-1] = time.Now()

						//finding the replace order from replace model and deleting it
						replaceOrderCollection := conn.ConnectDB("replaceOrders")
						replaceOrderFilter := bson.M{"productId": productId, "orderId": orderId}
						_, err := replaceOrderCollection.DeleteOne(context.TODO(), replaceOrderFilter)
						if err != nil {
							conn.GetError(err, w)
							return
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, "Cancel Expires!!")
						return
					}
				} else if oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Return Confirmed" {
					//code there
					nn := len(oldOrder.DeliveredDates[i])
					if time.Now().Before(oldOrder.DeliveredDates[i][nn-1]) {
						//write the code below there to place
						oldOrder.Status[i] = append(oldOrder.Status[i], "Cancelled")
						oldOrder.DeliveredDates[i][nn-1] = time.Now()

						//finding the return order from return model and deleting it
						returnOrderCollection := conn.ConnectDB("returnOrders")
						returnOrderFilter := bson.M{"productId": productId, "orderId": orderId}
						_, err := returnOrderCollection.DeleteOne(context.TODO(), returnOrderFilter)
						if err != nil {
							conn.GetError(err, w)
							return
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, "Cancel Expires!!")
						return
					}
				} else {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "You can't Cancel an order after Delivered!!")
					return
				}
			}
		}

		orderUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "status", Value: oldOrder.Status},
					{Key: "deliveredDates", Value: oldOrder.DeliveredDates},
				},
			},
		}
		_, errr := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}

		json.NewEncoder(w).Encode(oldOrder)
	}
}

func ChangeDeliveryDate(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)
		ids := params["orderId"]
		orderId, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idx := params["productId"]
		productId, errrr := primitive.ObjectIDFromHex(idx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idxx := params["id"]
		customerId, errrr := primitive.ObjectIDFromHex(idxx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//taking the day to increase the delivery from the body of the postman
		data, err := ioutil.ReadAll(r.Body)
		asString := string(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		days := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &days)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//fetching the order from database
		var oldOrder Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&oldOrder)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}

		//fetching hte customer from customer database
		var customer Model.Customer
		customerCollection := conn.ConnectDB("customers")
		customerFilter := bson.M{"_id": customerId}
		errrrr:= customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
		if errrrr != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		for i := 0; i < len(oldOrder.ProductIds); i++ {
			if oldOrder.ProductIds[i] == productId {
				if customer.IsPrime {
					if oldOrder.Status[i][len(oldOrder.Status[i])-1] == "Order Confirmed" {
						if oldOrder.DeliveredDates[i][1].After(time.Now()) {
							//write the code below there to change the delivery date
							oldOrder.DeliveredDates[i][1] = oldOrder.DeliveredDates[i][1].AddDate(0,0,int(days["noOfDayToIncrease"].(float64)))

							//updating the deliveryorder
							var deliverOrder Model.DeliverOrder
							deliverOrderCollection := conn.ConnectDB("deliverOrders")
							deliverOrderFilter := bson.M{"productId": productId, "orderId": orderId}
							errr := deliverOrderCollection.FindOne(context.TODO(), deliverOrderFilter).Decode(&deliverOrder)
							if errr != nil {
								conn.GetError(errr, w)
								return
							}
							deliverOrder.ExpectedDeliveryDate = oldOrder.DeliveredDates[i][1]
							deliverOrderUpdate := bson.D{
								{
									Key: "$set", Value: bson.D{
										{Key: "expectedDeliveryDate", Value: deliverOrder.ExpectedDeliveryDate},
									},
								},
							}
							deliverOrderResult,err := deliverOrderCollection.UpdateOne(context.TODO(), deliverOrderFilter, deliverOrderUpdate)
							if err != nil {
								conn.GetError(err, w)
								return
							}
							json.NewEncoder(w).Encode(deliverOrderResult)
						} else {
							w.WriteHeader(http.StatusBadRequest)
							fmt.Fprintln(w, "You can't change the Delivery date after the expected delivery date!!")
							return
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, "You can't change the Delivery date after the order been delivered!!")
						return
					}
				} else {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "You are not a prime member, kindly buy prime membership to change the delivery date!!")
					return
				}
				break
			}
		}

		orderUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "deliveredDates", Value: oldOrder.DeliveredDates},
				},
			},
		}
		_, errrrrr := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
		if errrrrr != nil {
			conn.GetError(errrrrr, w)
			return
		}

		json.NewEncoder(w).Encode(oldOrder)
	}
}
