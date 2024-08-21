package deliveryPartner

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDeliverOrders(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var deliverOrder []Model.DeliverOrder

	collection := conn.ConnectDB("deliverOrders")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.DeliverOrder
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		deliverOrder = append(deliverOrder, u)
	}
	json.NewEncoder(w).Encode(deliverOrder)
}
func GetReplaceOrders(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var replaceOrder []Model.ReplaceOrder

	collection := conn.ConnectDB("replaceOrders")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.ReplaceOrder
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		replaceOrder = append(replaceOrder, u)
	}
	json.NewEncoder(w).Encode(replaceOrder)
}
func GetReturnOrders(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var returnOrder []Model.ReturnOrder

	collection := conn.ConnectDB("returnOrders")

	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.ReturnOrder
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		returnOrder = append(returnOrder, u)
	}
	json.NewEncoder(w).Encode(returnOrder)
}

func OrderDelivered(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the orderid, productid and customerid from the url
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
		// idxx := params["customerId"]
		// customerId, errrr := primitive.ObjectIDFromHex(idxx)
		// if errrr != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		//taking the perticular order from order model to update the status of it
		var order Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&order)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}
		
		for i:=0;i<len(order.ProductIds);i++{
			if order.ProductIds[i] == productId{
				//taking the product which we need to update the status and delivery date
				var product Model.Product
				productCollection := conn.ConnectDB("products")
				productFilter := bson.M{"_id": productId}
				err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
				if err != nil {
					conn.GetError(err, w)
					return
				}

				if order.Status[i][len(order.Status[i])-1] != "Order Confirmed"{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Order is already Delivered!!")
					return
				}else{
					order.Status[i] = append(order.Status[i], "Order Delivered")
					n := len(order.DeliveredDates[i])
					order.DeliveredDates[i][n-1] = time.Now()
					orderUpdate := bson.D{
						{
							Key: "$set", Value: bson.D{
								{Key: "status", Value: order.Status},
								{Key: "deliveredDates", Value: order.DeliveredDates},
							},
						},
					}
					orderResult,err := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
					if err != nil {
						conn.GetError(err, w)
						return
					}
					json.NewEncoder(w).Encode(orderResult)

					//now here add report of this product for seller
					sellerId := product.SellerId
					var report Model.Report
					reportCollection := conn.ConnectDB("reports")
					reportFilter := bson.M{"sellerId": sellerId}
					errr := reportCollection.FindOne(context.TODO(), reportFilter).Decode(&report)
					if errr != nil {
						conn.GetError(errr, w)
						return
					}
					flag := 1
					for i:=0;i<len(report.SoldItems);i++{
						if report.SoldItems[i].ProductId == productId{
							report.SoldItems[i].ProductPrice = append(report.SoldItems[i].ProductPrice, product.SellingPrice)
							report.SoldItems[i].DeliveryDate = append(report.SoldItems[i].DeliveryDate, time.Now())
							report.SoldItems[i].Quantity = append(report.SoldItems[i].Quantity, order.OrderQuantitys[i])

							flag=0
							break
						}
					}
					if flag==1{
						var newSoldItem Model.SoldItem
						newSoldItem.ProductId = productId
						newSoldItem.ProductName = product.ProductName
						newSoldItem.ProductImage = product.ProductImages
						newSoldItem.ProductPrice = append(newSoldItem.ProductPrice, product.SellingPrice)
						newSoldItem.DeliveryDate = append(newSoldItem.DeliveryDate, time.Now())
						newSoldItem.Quantity = append(newSoldItem.Quantity, order.OrderQuantitys[i])

						report.SoldItems = append(report.SoldItems, newSoldItem)
					}
					//code here to update seller report
					reportUpdate := bson.D{
						{
							Key: "$set", Value: bson.D{
								{Key: "soldItems", Value: report.SoldItems},
							},
						},
					}
					reportResult,err := reportCollection.UpdateOne(context.TODO(), reportFilter, reportUpdate)
					if err != nil {
						conn.GetError(err, w)
						return
					}
					json.NewEncoder(w).Encode(reportResult)
				}

				//deleting the deliverOrder from the deliverOrder model
				deliverOrderCollection := conn.ConnectDB("deliverOrders")
				deliverOrderFilter := bson.M{"productId": productId, "orderId": orderId}
				deliverOrderResult, errr := deliverOrderCollection.DeleteOne(context.TODO(), deliverOrderFilter)
				if errr != nil {
					conn.GetError(errr, w)
					return
				}
				json.NewEncoder(w).Encode(deliverOrderResult)
			}
		}
		json.NewEncoder(w).Encode(order)
	}
}

func OrderReplaced(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the orderid, productid and customerid from the url
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
		// idxx := params["customerId"]
		// customerId, errrr := primitive.ObjectIDFromHex(idxx)
		// if errrr != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		//taking the perticular order from order model to update the status of it
		var order Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&order)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}
		
		for i:=0;i<len(order.ProductIds);i++{
			if order.ProductIds[i] == productId{
				//taking the replaceOrder from replaceOrder model
				var replaceOrder Model.ReplaceOrder
				replaceOrderCollection := conn.ConnectDB("replaceOrders")
				replaceOrderFilter := bson.M{"productId": productId, "orderId": orderId}
				errr := replaceOrderCollection.FindOne(context.TODO(), replaceOrderFilter).Decode(&replaceOrder)
				if errr != nil {
					conn.GetError(errr, w)
					return
				}

				//taking the product which we need to update the status and delivery date
				var product Model.Product
				productCollection := conn.ConnectDB("products")
				productFilter := bson.M{"_id": productId}
				err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
				if err != nil {
					conn.GetError(err, w)
					return
				}

				//updating the product quantity and unitsold according the issue giving by the customer
				if !replaceOrder.IsDamage && replaceOrder.DontLikeDueToColorOrSize{
					product.Quantity += order.OrderQuantitys[i]
					productUpdate := bson.D{
						{
							Key: "$set", Value: bson.D{
								{Key: "quantity", Value: product.Quantity},
							},
						},
					}
					productResult,err := productCollection.UpdateOne(context.TODO(), productFilter, productUpdate)
					if err != nil {
						conn.GetError(err, w)
						return
					}
					json.NewEncoder(w).Encode(productResult)
				}
				
				if order.Status[i][len(order.Status[i])-1] != "Replace Confirmed"{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Replace Order is already Delivered!!")
					return
				}else{
					order.Status[i] = append(order.Status[i], "Replace Order Delivered")
					n := len(order.DeliveredDates[i])
					order.DeliveredDates[i][n-1] = time.Now()
					orderUpdate := bson.D{
						{
							Key: "$set", Value: bson.D{
								{Key: "status", Value: order.Status},
								{Key: "deliveredDates", Value: order.DeliveredDates},
							},
						},
					}
					orderResult,err := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
					if err != nil {
						conn.GetError(err, w)
						return
					}
					json.NewEncoder(w).Encode(orderResult)
				}

				//deleting the replaceOrder from the replaceOrder model
				replaceOrderResult, errr := replaceOrderCollection.DeleteOne(context.TODO(), replaceOrderFilter)
				if errr != nil {
					conn.GetError(errr, w)
					return
				}
				json.NewEncoder(w).Encode(replaceOrderResult)
			}
		}
		json.NewEncoder(w).Encode(order)
	}
}

func OrderReturned(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		//taking the orderid, productid and customerid from the url
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
		idxx := params["customerId"]
		customerId, errrr := primitive.ObjectIDFromHex(idxx)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//taking the perticular order from order model to update the status of it
		var order Model.Order
		orderCollection := conn.ConnectDB("orders")
		orderFilter := bson.M{"_id": orderId}
		errr := orderCollection.FindOne(context.TODO(), orderFilter).Decode(&order)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}
		
		for i:=0;i<len(order.ProductIds);i++{
			if order.ProductIds[i] == productId{
				//taking the returnOrder from returnOrder model
				var returnOrder Model.ReturnOrder
				returnOrderCollection := conn.ConnectDB("returnOrders")
				returnOrderFilter := bson.M{"productId": productId, "orderId": orderId}
				errr := returnOrderCollection.FindOne(context.TODO(), returnOrderFilter).Decode(&returnOrder)
				if errr != nil {
					conn.GetError(errr, w)
					return
				}

				//taking the product which we need to update the status and delivery date
				var product Model.Product
				productCollection := conn.ConnectDB("products")
				productFilter := bson.M{"_id": productId}
				err := productCollection.FindOne(context.TODO(), productFilter).Decode(&product)
				if err != nil {
					conn.GetError(err, w)
					return
				}

				//updating the product quantity and unitsold according the issue giving by the customer
				product.UnitsSold -= order.OrderQuantitys[i]
				productUpdate := bson.D{
					{
						Key: "$set", Value: bson.D{
							{Key: "unitsSold", Value: product.UnitsSold},
						},
					},
				}
				if !returnOrder.IsDamage && returnOrder.DontLikeDueToColorOrSize{
					product.Quantity += order.OrderQuantitys[i]
					productUpdate = bson.D{
						{
							Key: "$set", Value: bson.D{
								{Key: "quantity", Value: product.Quantity},
								{Key: "unitsSold", Value: product.UnitsSold},
							},
						},
					}
				}
				productResult,err := productCollection.UpdateOne(context.TODO(), productFilter, productUpdate)
				if err != nil {
					conn.GetError(err, w)
					return
				}
				json.NewEncoder(w).Encode(productResult)
				
				//updating the status and return date 
				if order.Status[i][len(order.Status[i])-1] != "Return Confirmed"{
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w,"Return Order is already Taken or Only one time you can return the product!!")
					return
				}else{
					order.Status[i] = append(order.Status[i], "Refunded")
					n := len(order.DeliveredDates[i])
					order.DeliveredDates[i][n-1] = time.Now()
					orderUpdate := bson.D{
						{
							Key: "$set", Value: bson.D{
								{Key: "status", Value: order.Status},
								{Key: "deliveredDates", Value: order.DeliveredDates},
							},
						},
					}
					orderResult,err := orderCollection.UpdateOne(context.TODO(), orderFilter, orderUpdate)
					if err != nil {
						conn.GetError(err, w)
						return
					}
					json.NewEncoder(w).Encode(orderResult)
				}

				//deleting the returnOrder from the returnOrder model
				returnOrderResult, errr := returnOrderCollection.DeleteOne(context.TODO(), returnOrderFilter)
				if errr != nil {
					conn.GetError(errr, w)
					return
				}
				json.NewEncoder(w).Encode(returnOrderResult)

				//updating the customer wallet
				var customer Model.Customer
				cutomerCollection := conn.ConnectDB("customers")
				customerFilter := bson.M{"_id": customerId}
				errrr := cutomerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
				if errrr != nil {
					conn.GetError(errrr, w)
					return
				}
				customer.Wallet += product.SellingPrice * order.OrderQuantitys[i]
				customerUpdate := bson.D{
					{
						Key: "$set", Value: bson.D{
							{Key: "wallet", Value: customer.Wallet},
						},
					},
				}
				customerResult,err := cutomerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
				if err != nil {
					conn.GetError(err, w)
					return
				}
				json.NewEncoder(w).Encode(customerResult)
			}
		}
		json.NewEncoder(w).Encode(order)
	}
}