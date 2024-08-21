package customerCare

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	"strconv"

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

func GetCustomerCareMessages(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var customerCare []Model.CustomerCare
	//finding customerCare in database
	collection := conn.ConnectDB("customerCares")
	curr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.CustomerCare
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		customerCare = append(customerCare, u)
	}
	json.NewEncoder(w).Encode(customerCare)
}

func ChatWithMe(w http.ResponseWriter, r *http.Request) {
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
		message := make(map[string]interface{})
		err = json.Unmarshal([]byte(asString), &message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//fetching all the orders of the customer
		var yourOrders []Model.Order
		orderConnection := conn.ConnectDB("orders")
		orderFilter := bson.M{"customerId": customerId}
		orderSort := options.Find()
		orderSort.SetSort(bson.D{{"orderedDate", -1}})
		curr, err := orderConnection.Find(context.TODO(), orderFilter, orderSort)
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
			yourOrders = append(yourOrders, u)
		}

		//messaging start here
		if message["message"].(string) == "" {
			m := "Hi! It's Amazon's messaging assistant again."
			fmt.Fprintln(w, m)
			fmt.Fprintln(w, "Product Name: "+yourOrders[0].ProductNames[0])
			fmt.Fprintln(w, "Is this what you need help with?")

			//recommendations
			fmt.Fprintln(w, "Choose Anyone of the following...")
			fmt.Fprintln(w, "No, Something else")
			fmt.Fprintln(w, "Yes, that's correct")
		}

		if message["message"].(string) == "Yes, that's correct" {
			status:= yourOrders[0].Status[0][0]
			fmt.Fprintln(w,"Status for the Product: "+status+".")
			if status == "Refunded"{
				year := strconv.Itoa(yourOrders[0].DeliveredDates[0][len(yourOrders[0].DeliveredDates[0])-1].Year())
				month := strconv.Itoa(int(yourOrders[0].DeliveredDates[0][len(yourOrders[0].DeliveredDates[0])-1].Month()))
				day := strconv.Itoa(yourOrders[0].DeliveredDates[0][len(yourOrders[0].DeliveredDates[0])-1].Day())
				fmt.Fprintln(w,"We started the refund process on "+day+"/"+month+"/"+year+". The refund will get added to your original payment method within 5-7 business days.")
				fmt.Fprintln(w,"In case your refund is delayed, please contact your bank.")
			}
			//recommendation
			fmt.Fprintln(w, "Choose Anyone of the following...")
			fmt.Fprintln(w, "Ok, Thanks!!")
		} else if message["message"].(string) == "No, Something else" {
			fmt.Fprintln(w, "So, what can I help you with?")

			//recommendations
			fmt.Fprintln(w, "Choose Anyone of the following...")
			fmt.Fprintln(w, "An item I ordered")
			fmt.Fprintln(w, "Something else")
		}

		if message["message"].(string) == "Ok, Thanks!!"{
			fmt.Fprintln(w,"Thanks for sharing.")
		}

		var productName []string
		if message["message"].(string) == "An item I ordered" {
			cnt:=0;
			flag:=1
			fmt.Fprintln(w,"Lets see. Could you select the item you're looking for from your recent orders below?")
			for i:=0;i<len(yourOrders);i++{
				for j:=0 ;j<len(yourOrders[i].ProductIds);j++{
					fmt.Fprintln(w, "Product Name: "+yourOrders[i].ProductNames[j])
					// fmt.Fprintln(w, "Ordered Date: "+strconv.Itoa(yourOrders[i].OrderedDate.(string)))
					productName=append(productName, yourOrders[i].ProductNames[j])
					cnt++
					if cnt>=5{
						flag=0
						break
					}
				}
				if flag==0{
					break
				}
			}
		}else if message["message"].(string) == "Something else" {
			fmt.Fprintln(w, "Request for phone call!!")
			fmt.Fprintln(w,"Please wait, We are calling....")
		}

		for i:=0;i<len(productName);i++{
			if productName[i]==message["message"].(string){
				fmt.Fprintln(w,"Ok, you can request a call back.")

				//recommendations
				fmt.Fprintln(w, "Choose Anyone of the following...")
				fmt.Fprintln(w, "Request a phone call with agent")
				fmt.Fprintln(w, "I don't need more help")
			}
		}

		if message["message"].(string) == "Request a phone call with agent"{
			fmt.Fprintln(w,"Please wait, we are calling....")
		}else if message["message"].(string) == "I don't need more help"{
			fmt.Fprintln(w,"Thanks for sharing.")
		}
		// var customerCare Model.CustomerCare
		// newcustomerCareId := primitive.NewObjectID()
		// customerCare.Id = newcustomerCareId
		// collection := conn.ConnectDB("customerCares")
		// result, err := collection.InsertOne(context.TODO(), customerCare)
		// if err != nil {
		// 	conn.GetError(err, w)
		// 	return
		// }

		// customerCollection := conn.ConnectDB("customers")
		// filter := bson.M{"_id": customerId}
		// errr := customerCollection.FindOne(context.TODO(), filter).Decode(&customer)
		// if errr != nil {
		// 	conn.GetError(err, w)
		// 	return
		// }
		// update := bson.D{
		// 	{
		// 		Key: "$set", Value: bson.D{
		// 			{Key: "customerCares", Value: customer.customerCares},
		// 		},
		// 	},
		// }
		// customerResult, errrr := customerCollection.UpdateOne(context.TODO(), filter, update)
		// if errrr != nil {
		// 	conn.GetError(err, w)
		// 	return
		// }
		// json.NewEncoder(w).Encode(result)
		// json.NewEncoder(w).Encode(customerResult)
	}
}

// func DeleteCustomerCare(w http.ResponseWriter, r *http.Request) {
// 	Generic.SetupResponse(&w, r)
// 	if r.Method == "DELETE" {
// 		w.Header().Set("Content-Type", "application/json")

// 		var params = mux.Vars(r)

// 		ids := params["customerCareId"]
// 		id, _ := primitive.ObjectIDFromHex(ids)

// 		filter := bson.M{"_id": id}
// 		collection := conn.ConnectDB("customerCares")

// 		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
// 		if err != nil {
// 			conn.GetError(err, w)
// 			return
// 		}

// 		//deleting the customerCare in the customer model
// 		idd := params["id"]
// 		customerId, errrr := primitive.ObjectIDFromHex(idd)
// 		if errrr != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		var customer Model.Customer
// 		customerCollection := conn.ConnectDB("customers")
// 		customerFilter := bson.M{"_id": customerId}
// 		errr := customerCollection.FindOne(context.TODO(), customerFilter).Decode(&customer)
// 		if errr != nil {
// 			conn.GetError(err, w)
// 			return
// 		}
// 		customerUpdate := bson.D{
// 			{
// 				Key: "$set", Value: bson.D{
// 					{Key: "customerCares", Value: customer.customerCares},
// 				},
// 			},
// 		}
// 		customerResult, errrr := customerCollection.UpdateOne(context.TODO(), customerFilter, customerUpdate)
// 		if errrr != nil {
// 			conn.GetError(err, w)
// 			return
// 		}
// 		json.NewEncoder(w).Encode(deleteResult)
// 		json.NewEncoder(w).Encode(customerResult)
// 	}

// }
