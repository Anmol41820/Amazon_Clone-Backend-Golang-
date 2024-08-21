package product

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	// "io"
	"math"
	// "os"
	"strconv"
	// "sync"
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
// var storageMutex sync.Mutex

func GetProduct(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	var product []Model.Product

	//adding pagination in products
	limit := 2
	var params = mux.Vars(r)
	pageNumber, ok := params["pageNumber"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(pageNumber)
	pageFilter := options.Find()
	pageFilter.SetSkip(int64(page-1) * int64(limit))
	pageFilter.SetLimit(int64(limit))

	//finding product in database
	collection := conn.ConnectDB("products")
	curr, err := collection.Find(context.TODO(), bson.M{}, pageFilter)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var u Model.Product
		err := curr.Decode(&u)
		if err != nil {
			fmt.Println("****ERROR*****")
			w.WriteHeader(http.StatusBadGateway)
		}
		product = append(product, u)
	}

	if err := curr.Err(); err != nil {
		// log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(product)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
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
		id, errrr := primitive.ObjectIDFromHex(ids)
		if errrr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product := make(map[string]interface{})
		product["averageRating"] = 0.0
		product["unitsSold"] = 0
		product["sellerId"] = id
		product["releaseDate"] = time.Now()
		product["newRelease"] = true
		product["bestSeller"] = false

		err = json.Unmarshal([]byte(asString), &product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if product == nil {
			http.Error(w, "Invalid product data", http.StatusBadRequest)
			return
		}
		productId := primitive.NewObjectID()
		product["_id"] = productId
		f := product["maxRetailPrice"].(float64)
		i := product["sellingPrice"].(float64)
		dis := ((f - i) / f) * 100
		dis = math.Round(dis*100) / 100
		product["discount"] = dis

		//dealing with product images
		// files := r.MultipartForm.File["productImages"]
		// var imagePaths []string
		// for _, fileHeader := range files {
		// 	// Open the file
		// 	file, err := fileHeader.Open()
		// 	if err != nil {
		// 		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		// 		return
		// 	}
		// 	defer file.Close()

		// 	// Save the image locally
		// 	storageMutex.Lock()
		// 	defer storageMutex.Unlock()
		// 	// Create the file in the local storage directory
		// 	productName := product["productName"].(string)
		// 	fileName := fileHeader.Filename
		// 	outFile, errr := os.Create("/ProductImages/"+productName+"/"+fileName)
		// 	if errr != nil {
		// 		fmt.Println("first")
		// 		http.Error(w, errr.Error(), http.StatusBadRequest)
		// 		return
		// 	}
		// 	defer outFile.Close()
		// 	// Copy the image data to the file
		// 	_, err = io.Copy(outFile, file)
		// 	if err != nil {
		// 		fmt.Println("second")
		// 		http.Error(w, err.Error(), http.StatusBadRequest)
		// 		return
		// 	}

		// 	// Add the file path to the array
		// 	imagePaths = append(imagePaths, fileHeader.Filename)
		// }
		// product["productImages"] = imagePaths

		collection := conn.ConnectDB("products")
		productResult, err := collection.InsertOne(context.TODO(), product)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		//inserting productid in seller model
		var seller Model.Seller
		sellerCollection := conn.ConnectDB("sellers")
		sellerFilter := bson.M{"_id": id}
		errr := sellerCollection.FindOne(context.TODO(), sellerFilter).Decode(&seller)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}
		seller.ProductsListedIds = append(seller.ProductsListedIds, productId)
		sellerUpdate := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "productsListedIds", Value: seller.ProductsListedIds},
				},
			},
		}
		sellerResult, errrrr := sellerCollection.UpdateOne(context.TODO(), sellerFilter, sellerUpdate)
		if errrrr != nil {
			conn.GetError(errrrr, w)
			return
		}
		json.NewEncoder(w).Encode(productResult)
		json.NewEncoder(w).Encode(sellerResult)
	}
}

func GetSingleProduct(w http.ResponseWriter, r *http.Request) {

	Generic.SetupResponse(&w, r)

	w.Header().Set("Content-Type", "application/json")
	productCollection := conn.ConnectDB("products")

	var product Model.Product
	var params = mux.Vars(r)
	ids := params["productId"]
	id, _ := primitive.ObjectIDFromHex(ids)
	filter := bson.M{"_id": id}
	err := productCollection.FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		conn.GetError(err, w)
		return
	}

	//fetching top 10 product based on category of the selected product
	var recommendatedProducts []Model.Product
	for j := 0; j < len(product.ProductCategories); j++ {
		recommendationFilter := bson.M{
			"productCategories": bson.M{"$in": []string{product.ProductCategories[j]}},
		}
		limitFilter := options.Find()
		limitFilter.SetLimit(10)
		curr, errr := productCollection.Find(context.TODO(), recommendationFilter, limitFilter)
		if errr != nil {
			conn.GetError(errr, w)
			return
		}

		for curr.Next(context.TODO()) {
			var u Model.Product
			err := curr.Decode(&u)
			if err != nil {
				conn.GetError(err, w)
				return
			}
			recommendatedProducts = append(recommendatedProducts, u)
			if len(recommendatedProducts) >= 10 {
				break
			}
		}
		if len(recommendatedProducts) >= 10 {
			break
		}
	}

	//updation the productRecommendations API
	idx := params["id"]
	customerId, _ := primitive.ObjectIDFromHex(idx)
	var customerProductRecommendation Model.ProductRecommendation

	customerProductRecommendationCollection := conn.ConnectDB("productRecommendations")
	customerProductRecommendationFilter := bson.M{"customerId": customerId}
	errr := customerProductRecommendationCollection.FindOne(context.TODO(), customerProductRecommendationFilter).Decode(&customerProductRecommendation)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}
	customerProductRecommendation.ProductIds = customerProductRecommendation.ProductIds[:0]
	customerProductRecommendation.AboutProduct = customerProductRecommendation.AboutProduct[:0]
	customerProductRecommendation.AverageRating = customerProductRecommendation.AverageRating[:0]
	customerProductRecommendation.BestSeller = customerProductRecommendation.BestSeller[:0]
	customerProductRecommendation.Brand = customerProductRecommendation.Brand[:0]
	customerProductRecommendation.Discount = customerProductRecommendation.Discount[:0]
	customerProductRecommendation.MaxRetailPrice = customerProductRecommendation.MaxRetailPrice[:0]
	customerProductRecommendation.NewRelease = customerProductRecommendation.NewRelease[:0]
	customerProductRecommendation.NumberOfReviews = customerProductRecommendation.NumberOfReviews[:0]
	customerProductRecommendation.ProductImage = customerProductRecommendation.ProductImage[:0]
	customerProductRecommendation.ProductName = customerProductRecommendation.ProductName[:0]
	customerProductRecommendation.SellingPrice = customerProductRecommendation.SellingPrice[:0]
	for i := 0; i < len(recommendatedProducts); i++ {
		customerProductRecommendation.ProductIds = append(customerProductRecommendation.ProductIds, recommendatedProducts[i].Id)
		customerProductRecommendation.AboutProduct = append(customerProductRecommendation.AboutProduct, recommendatedProducts[i].AboutProduct)
		customerProductRecommendation.AverageRating = append(customerProductRecommendation.AverageRating, recommendatedProducts[i].AverageRating)
		customerProductRecommendation.BestSeller = append(customerProductRecommendation.BestSeller, recommendatedProducts[i].BestSeller)
		customerProductRecommendation.Brand = append(customerProductRecommendation.Brand, recommendatedProducts[i].Brand)
		customerProductRecommendation.Discount = append(customerProductRecommendation.Discount, recommendatedProducts[i].Discount)
		customerProductRecommendation.MaxRetailPrice = append(customerProductRecommendation.MaxRetailPrice, recommendatedProducts[i].MaxRetailPrice)
		customerProductRecommendation.NewRelease = append(customerProductRecommendation.NewRelease, recommendatedProducts[i].NewRelease)
		customerProductRecommendation.NumberOfReviews = append(customerProductRecommendation.NumberOfReviews, recommendatedProducts[i].NumberOfReviews)
		customerProductRecommendation.ProductImage = append(customerProductRecommendation.ProductImage, recommendatedProducts[i].ProductImages[0])
		customerProductRecommendation.ProductName = append(customerProductRecommendation.ProductName, recommendatedProducts[i].ProductName)
		customerProductRecommendation.SellingPrice = append(customerProductRecommendation.SellingPrice, recommendatedProducts[i].SellingPrice)
	}
	customerProductRecommendationUpdate := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "productIds", Value: customerProductRecommendation.ProductIds},
				{Key: "productName", Value: customerProductRecommendation.ProductName},
				{Key: "aboutProduct", Value: customerProductRecommendation.AboutProduct},
				{Key: "brand", Value: customerProductRecommendation.Brand},
				{Key: "bestSeller", Value: customerProductRecommendation.BestSeller},
				{Key: "newRelease", Value: customerProductRecommendation.NewRelease},
				{Key: "maxRetailPrice", Value: customerProductRecommendation.MaxRetailPrice},
				{Key: "sellingPrice", Value: customerProductRecommendation.SellingPrice},
				{Key: "discount", Value: customerProductRecommendation.Discount},
				{Key: "averageRating", Value: customerProductRecommendation.AverageRating},
				{Key: "productImage", Value: customerProductRecommendation.ProductImage},
				{Key: "numberOfReviews", Value: customerProductRecommendation.NumberOfReviews},
			},
		},
	}
	customerProductRecommendationResult, errrrr := customerProductRecommendationCollection.UpdateOne(context.TODO(), customerProductRecommendationFilter, customerProductRecommendationUpdate)
	if errrrr != nil {
		conn.GetError(errrrr, w)
		return
	}

	//updating the recentlyViewedProduct of the customer
	var recentlyViewedProduct Model.RecentlyViewedProduct
	recentlyViewedProductCollection := conn.ConnectDB("recentlyViewedProducts")
	recentlyViewedProductFilter := bson.M{"customerId": customerId}
	errrr := recentlyViewedProductCollection.FindOne(context.TODO(), recentlyViewedProductFilter).Decode(&recentlyViewedProduct)
	if errrr != nil {
		conn.GetError(errrr, w)
		return
	}
	n := len(recentlyViewedProduct.ProductIds)
	flag := 1
	for k := 0; k < n; k++ {
		if recentlyViewedProduct.ProductIds[k] == product.Id {
			flag = 0
			break
		}
	}
	if flag == 1 {
		recentlyViewedProduct.ProductIds = append(recentlyViewedProduct.ProductIds, product.Id)
		recentlyViewedProduct.AboutProduct = append(recentlyViewedProduct.AboutProduct, product.AboutProduct)
		recentlyViewedProduct.AverageRating = append(recentlyViewedProduct.AverageRating, product.AverageRating)
		recentlyViewedProduct.BestSeller = append(recentlyViewedProduct.BestSeller, product.BestSeller)
		recentlyViewedProduct.Brand = append(recentlyViewedProduct.Brand, product.Brand)
		recentlyViewedProduct.Discount = append(recentlyViewedProduct.Discount, product.Discount)
		recentlyViewedProduct.MaxRetailPrice = append(recentlyViewedProduct.MaxRetailPrice, product.MaxRetailPrice)
		recentlyViewedProduct.NewRelease = append(recentlyViewedProduct.NewRelease, product.NewRelease)
		recentlyViewedProduct.NumberOfReviews = append(recentlyViewedProduct.NumberOfReviews, product.NumberOfReviews)
		recentlyViewedProduct.ProductImage = append(recentlyViewedProduct.ProductImage, product.ProductImages[0])
		recentlyViewedProduct.ProductName = append(recentlyViewedProduct.ProductName, product.ProductName)
		recentlyViewedProduct.SellingPrice = append(recentlyViewedProduct.SellingPrice, product.SellingPrice)
	}

	if n >= 10 {
		recentlyViewedProduct.ProductIds = recentlyViewedProduct.ProductIds[n-10:]
		recentlyViewedProduct.ProductName = recentlyViewedProduct.ProductName[n-10:]
		recentlyViewedProduct.AboutProduct = recentlyViewedProduct.AboutProduct[n-10:]
		recentlyViewedProduct.Brand = recentlyViewedProduct.Brand[n-10:]
		recentlyViewedProduct.BestSeller = recentlyViewedProduct.BestSeller[n-10:]
		recentlyViewedProduct.NewRelease = recentlyViewedProduct.NewRelease[n-10:]
		recentlyViewedProduct.MaxRetailPrice = recentlyViewedProduct.MaxRetailPrice[n-10:]
		recentlyViewedProduct.SellingPrice = recentlyViewedProduct.SellingPrice[n-10:]
		recentlyViewedProduct.Discount = recentlyViewedProduct.Discount[n-10:]
		recentlyViewedProduct.AverageRating = recentlyViewedProduct.AverageRating[n-10:]
		recentlyViewedProduct.ProductImage = recentlyViewedProduct.ProductImage[n-10:]
		recentlyViewedProduct.NumberOfReviews = recentlyViewedProduct.NumberOfReviews[n-10:]
	}
	recentlyViewedProductUpdate := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "productIds", Value: recentlyViewedProduct.ProductIds},
				{Key: "productName", Value: recentlyViewedProduct.ProductName},
				{Key: "aboutProduct", Value: recentlyViewedProduct.AboutProduct},
				{Key: "brand", Value: recentlyViewedProduct.Brand},
				{Key: "bestSeller", Value: recentlyViewedProduct.BestSeller},
				{Key: "newRelease", Value: recentlyViewedProduct.NewRelease},
				{Key: "maxRetailPrice", Value: recentlyViewedProduct.MaxRetailPrice},
				{Key: "sellingPrice", Value: recentlyViewedProduct.SellingPrice},
				{Key: "discount", Value: recentlyViewedProduct.Discount},
				{Key: "averageRating", Value: recentlyViewedProduct.AverageRating},
				{Key: "productImage", Value: recentlyViewedProduct.ProductImage},
				{Key: "numberOfReviews", Value: recentlyViewedProduct.NumberOfReviews},
			},
		},
	}
	recentlyViewedProductResult, errrrr := recentlyViewedProductCollection.UpdateOne(context.TODO(), recentlyViewedProductFilter, recentlyViewedProductUpdate)
	if errrrr != nil {
		conn.GetError(errrrr, w)
		return
	}

	json.NewEncoder(w).Encode(customerProductRecommendationResult)
	json.NewEncoder(w).Encode(recentlyViewedProductResult)
	json.NewEncoder(w).Encode(product)
	json.NewEncoder(w).Encode(recommendatedProducts)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	// if(!conn.ProtectedHandler(w,r)){return}

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["productId"]
		id, _ := primitive.ObjectIDFromHex(ids)
		var product Model.Product

		filter := bson.M{"_id": id}

		_ = json.NewDecoder(r.Body).Decode(&product)

		f := float64(product.MaxRetailPrice)
		i := float64(product.SellingPrice)
		dis := ((f - i) / f) * 100
		dis = math.Round(dis*100) / 100
		product.Discount = dis

		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "aboutProduct", Value: product.AboutProduct},
					{Key: "brand", Value: product.Brand},
					{Key: "color", Value: product.Color},
					{Key: "discount", Value: product.Discount},
					{Key: "maxRetailPrice", Value: product.MaxRetailPrice},
					{Key: "productCategories", Value: product.ProductCategories},
					{Key: "productImages", Value: product.ProductImages},
					{Key: "productName", Value: product.ProductName},
					{Key: "productProperties", Value: product.ProductProperties},
					{Key: "quantity", Value: product.Quantity},
					// {Key: "reviews", Value: product.Reviews},
					{Key: "sellingPrice", Value: product.SellingPrice},
					// {Key: "unitsSold", Value: product.UnitsSold},
				},
			},
		}

		collection := conn.ConnectDB("products")
		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(product)
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")

		var params = mux.Vars(r)

		ids := params["productId"]
		id, _ := primitive.ObjectIDFromHex(ids)

		filter := bson.M{"_id": id}
		collection := conn.ConnectDB("products")
		deleteResult, err := collection.DeleteOne(context.TODO(), filter)
		if err != nil {
			conn.GetError(err, w)
			return
		}

		json.NewEncoder(w).Encode(deleteResult)
	}

}

func SearchProduct(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	asString := string(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	search := make(map[string]string)
	err = json.Unmarshal([]byte(asString), &search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if search == nil {
		http.Error(w, "Invalid search data", http.StatusBadRequest)
		return
	}
	// delete(search, "_id")

	//adding pagination in search
	limit := 2
	var params = mux.Vars(r)
	pageNumber, ok := params["pageNumber"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(pageNumber)
	pageFilter := options.Find()
	pageFilter.SetSkip(int64(page-1) * int64(limit))
	pageFilter.SetLimit(int64(limit))

	//finding product according the search
	collection := conn.ConnectDB("products")
	filter := bson.M{"productName": bson.M{"$regex": search["search_str"], "$options": "i"}}
	curr, errr := collection.Find(context.TODO(), filter, pageFilter)
	if errr != nil {
		conn.GetError(err, w)
		return
	}
	var products []Model.Product
	for curr.Next(context.TODO()) {
		var u Model.Product
		err := curr.Decode(&u)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		products = append(products, u)
	}
	json.NewEncoder(w).Encode(products)

	//taking the customerId form url
	idxx := params["id"]
	customerId, errrr := primitive.ObjectIDFromHex(idxx)
	if errrr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//inserting the search_str in search history model
	var searchHistory Model.SearchHistory
	searchHistoryCollection := conn.ConnectDB("searchHistories")
	searchHistoryFilter := bson.M{"customerId": customerId}
	errrrr := searchHistoryCollection.FindOne(context.TODO(), searchHistoryFilter).Decode(&searchHistory)
	if errrrr != nil {
		conn.GetError(errrrr, w)
		return
	}
	if len(searchHistory.SearchText) < 10 {
		searchHistory.SearchText = append(searchHistory.SearchText, search["search_str"])
	} else {
		searchHistory.SearchText = append(searchHistory.SearchText, search["search_str"])
		searchHistory.SearchText = searchHistory.SearchText[len(searchHistory.SearchText)-10:]
	}
	searchHistoryUpdate := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "searchText", Value: searchHistory.SearchText},
			},
		},
	}
	searchHistoryResult, err := searchHistoryCollection.UpdateOne(context.TODO(), searchHistoryFilter, searchHistoryUpdate)
	if err != nil {
		conn.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(searchHistoryResult)
}

func CategoryFilter(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(r.Body)
	asString := string(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoryFilter := make(map[string]interface{})
	categoryFilter["color"] = ""
	categoryFilter["categoryName"] = ""
	categoryFilter["brand"] = ""
	categoryFilter["minPrice"] = -1
	err = json.Unmarshal([]byte(asString), &categoryFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if categoryFilter == nil {
		http.Error(w, "Invalid categoryFilter data", http.StatusBadRequest)
		return
	}
	// delete(search, "_id")

	//adding pagination in search
	limit := 2
	var params = mux.Vars(r)
	pageNumber, ok := params["pageNumber"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(pageNumber)
	pageFilter := options.Find()
	pageFilter.SetSkip(int64(page-1) * int64(limit))
	pageFilter.SetLimit(int64(limit))

	//finding product according the search
	collection := conn.ConnectDB("products")
	filter := bson.M{
		"productCategories": bson.M{"$in": []string{categoryFilter["categoryName"].(string)}},
	}
	if categoryFilter["brand"] != "" && categoryFilter["minPrice"] == -1 {
		filter = bson.M{
			"productCategories": bson.M{"$in": []string{categoryFilter["categoryName"].(string)}},
			"brand":             categoryFilter["brand"].(string),
		}
	}
	if categoryFilter["minPrice"] != -1 && categoryFilter["brand"] != "" {
		filter = bson.M{
			"productCategories": bson.M{"$in": []string{categoryFilter["categoryName"].(string)}},
			"brand":             categoryFilter["brand"].(string),
			"sellingPrice": bson.M{
				"$gte": categoryFilter["minPrice"],
				"$lte": categoryFilter["maxPrice"],
			},
		}
	}
	if categoryFilter["minPrice"] != -1 && categoryFilter["brand"] == "" {
		filter = bson.M{
			"productCategories": bson.M{"$in": []string{categoryFilter["categoryName"].(string)}},
			"sellingPrice": bson.M{
				"$gte": categoryFilter["minPrice"],
				"$lte": categoryFilter["maxPrice"],
			},
		}
	}
	if categoryFilter["color"] != "" {
		filter = bson.M{
			"productCategories": bson.M{"$in": []string{categoryFilter["categoryName"].(string)}},
			"brand":             categoryFilter["brand"].(string),
			"color":             categoryFilter["color"].(string),
			"sellingPrice": bson.M{
				"$gte": categoryFilter["minPrice"],
				"$lte": categoryFilter["maxPrice"],
			},
		}
	}
	sortFilter := options.Find()
	if categoryFilter["sortPriceAcending"] == true {
		sortFilter.SetSort(bson.D{{"sellingPrice", 1}})
	}
	if categoryFilter["sortPriceDecending"] == true {
		sortFilter.SetSort(bson.D{{"sellingPrice", -1}})
	}
	curr, errr := collection.Find(context.TODO(), filter, sortFilter, pageFilter)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}
	var products []Model.Product
	for curr.Next(context.TODO()) {
		var u Model.Product
		err := curr.Decode(&u)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		products = append(products, u)
	}
	json.NewEncoder(w).Encode(products)
}

func CategoryFilterByBestSeller(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//adding pagination in search
	limit := 2
	var params = mux.Vars(r)
	pageNumber, ok := params["pageNumber"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(pageNumber)
	pageFilter := options.Find()
	pageFilter.SetSkip(int64(page-1) * int64(limit))
	pageFilter.SetLimit(int64(limit))

	//finding product according the best seller
	collection := conn.ConnectDB("products")
	filter := bson.M{
		"bestSeller": true,
	}
	sortFilter := options.Find()
	sortFilter.SetSort(bson.D{{"unitsSold", -1}})

	curr, errr := collection.Find(context.TODO(), filter, sortFilter, pageFilter)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}
	var products []Model.Product
	for curr.Next(context.TODO()) {
		var u Model.Product
		err := curr.Decode(&u)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		products = append(products, u)
	}
	json.NewEncoder(w).Encode(products)
}

func CategoryFilterByNewRelease(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//adding pagination in search
	limit := 2
	var params = mux.Vars(r)
	pageNumber, ok := params["pageNumber"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(pageNumber)
	pageFilter := options.Find()
	pageFilter.SetSkip(int64(page-1) * int64(limit))
	pageFilter.SetLimit(int64(limit))

	//finding product according the best seller
	collection := conn.ConnectDB("products")
	filter := bson.M{
		"newRelease": true,
	}
	sortFilter := options.Find()
	sortFilter.SetSort(bson.D{{"releaseDate", -1}})

	curr, errr := collection.Find(context.TODO(), filter, sortFilter, pageFilter)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}
	var products []Model.Product
	for curr.Next(context.TODO()) {
		var u Model.Product
		err := curr.Decode(&u)
		if err != nil {
			conn.GetError(err, w)
			return
		}
		products = append(products, u)
	}
	json.NewEncoder(w).Encode(products)
}
