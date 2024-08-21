package router

import (
	HomeController "Amazon_Server/Middleware/Home"
	UserController "Amazon_Server/Middleware/User"
	RegisterController "Amazon_Server/Middleware/Register"
	LoginController "Amazon_Server/Middleware/Login" 
	LogoutController "Amazon_Server/Middleware/Logout"
	ForgotPasswordController "Amazon_Server/Middleware/ForgotPassword"
	changePasswordController "Amazon_Server/Middleware/ChangePassword"
	CustomerController "Amazon_Server/Middleware/Customer"
	SellerController "Amazon_Server/Middleware/Seller"
	ProductController "Amazon_Server/Middleware/Product"
	CategoryController "Amazon_Server/Middleware/Category"
	CartController "Amazon_Server/Middleware/Cart"
	WishlistController "Amazon_Server/Middleware/Wishlist"
	AddressController "Amazon_Server/Middleware/Address"
	OrderController "Amazon_Server/Middleware/Order"
	DeliveryPartnerController "Amazon_Server/Middleware/DeliveryPartner"
	ReviewController "Amazon_Server/Middleware/Review"
	ReportController "Amazon_Server/Middleware/Report"
	SearchHistoryController "Amazon_Server/Middleware/SearchHistory"
	CustomerCareController "Amazon_Server/Middleware/CustomerCare"
	ProductRecommendationController "Amazon_Server/Middleware/ProductRecommendation"
	RecentlyViewedProductController "Amazon_Server/Middleware/RecentlyViewedProduct"
	

	auth "Amazon_Server/Middleware/Authorization"
	"net/http"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	apiRouter := mux.NewRouter()
	homeHandler(apiRouter)
	userHandler(apiRouter)
	customerHandler(apiRouter)
	sellerHandler(apiRouter)
	productHandler(apiRouter)
	categoryHandler(apiRouter)
	cartHandler(apiRouter)
	wishlistHandler(apiRouter)
	addressHandler(apiRouter)
	orderHandler(apiRouter)
	deliveryPartnerHandler(apiRouter)
	reviewHandler(apiRouter)
	reportHandler(apiRouter)
	searchHistoryHandler(apiRouter)
	customerCareHandler(apiRouter)
	productRecommendationHandler(apiRouter)
	recentlyViewedProductHandler(apiRouter)
	return apiRouter
}

func homeHandler(router *mux.Router){
	router.HandleFunc("/customerId={id}", HomeController.CheckPrimeMemberShip).Methods("GET")
}

func userHandler(router *mux.Router) {
	router.HandleFunc("/users", UserController.GetUser).Methods("GET")
	router.HandleFunc("/user/userId={id}", UserController.GetSingleUser).Methods("GET")
	router.HandleFunc("/addUser", UserController.CreateUser).Methods("POST", "OPTIONS")
	router.Handle("/updateUser/userId={id}", auth.ProtectedMiddleware(http.HandlerFunc(UserController.UpdateUser))).Methods("PUT", "OPTIONS")
	router.Handle("/deleteUser/userId={id}", auth.ProtectedMiddleware(http.HandlerFunc(UserController.DeleteUser))).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/registerCustomer", RegisterController.RegisterCustomer).Methods("POST", "OPTIONS")
	router.HandleFunc("/registerSeller", RegisterController.RegisterSeller).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", LoginController.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/logout", LogoutController.Logout).Methods("POST", "OPTIONS")
	router.HandleFunc("/forgotPassword", ForgotPasswordController.ForgotPassword).Methods("POST", "OPTIONS")
	router.HandleFunc("/changePassword", changePasswordController.ChangePassword).Methods("POST", "OPTIONS")
}
func customerHandler(router *mux.Router) {
	router.HandleFunc("/customers", CustomerController.GetCustomer).Methods("GET")
	router.HandleFunc("/customer/customerId={id}", CustomerController.GetSingleCustomer).Methods("GET")
	router.HandleFunc("/addCustomer", CustomerController.CreateCustomer).Methods("POST", "OPTIONS")
	router.Handle("/updateCustomer/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CustomerController.UpdateCustomer))).Methods("PUT", "OPTIONS")
	router.Handle("/deleteCustomer/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CustomerController.DeleteCustomer))).Methods("DELETE", "OPTIONS")

	router.Handle("/addMoneyInWallet/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CustomerController.AddMoneyInWallet))).Methods("POST", "OPTIONS")
	router.Handle("/withdrawMoneyFromWallet/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CustomerController.WithdrawMoneyFromWallet))).Methods("POST", "OPTIONS")
	router.Handle("/buyPrimeMembership/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CustomerController.BuyPrimeMembership))).Methods("POST", "OPTIONS")
}
func sellerHandler(router *mux.Router) {
	router.HandleFunc("/sellers", SellerController.GetSeller).Methods("GET")
	router.HandleFunc("/seller/sellerId={id}", SellerController.GetSingleSeller).Methods("GET")
	router.HandleFunc("/addSeller", SellerController.CreateSeller).Methods("POST", "OPTIONS")
	router.Handle("/updateSeller/sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(SellerController.UpdateSeller))).Methods("PUT", "OPTIONS")
	router.Handle("/deleteSeller/sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(SellerController.DeleteSeller))).Methods("DELETE", "OPTIONS")
}
func productHandler(router *mux.Router) {
	router.HandleFunc("/products/page={pageNumber}", ProductController.GetProduct).Methods("GET")
	// router.HandleFunc("/product/productId={productId}", ProductController.GetSingleProduct).Methods("GET")
	router.Handle("/product/productId={productId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ProductController.GetSingleProduct))).Methods("GET")
	router.Handle("/addProduct/sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ProductController.CreateProduct))).Methods("POST", "OPTIONS")
	router.Handle("/updateProduct/productId={productId}&sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ProductController.UpdateProduct))).Methods("PUT", "OPTIONS")
	router.Handle("/deleteProduct/productId={productId}&sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ProductController.DeleteProduct))).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/products/search/page={pageNumber}&customerId={id}", ProductController.SearchProduct).Methods("GET")
	router.HandleFunc("/products/categoryFilter/page={pageNumber}", ProductController.CategoryFilter).Methods("GET")
	router.HandleFunc("/products/categoryFilter/bestSeller/page={pageNumber}", ProductController.CategoryFilterByBestSeller).Methods("GET")
	router.HandleFunc("/products/categoryFilter/newRelease/page={pageNumber}", ProductController.CategoryFilterByNewRelease).Methods("GET")
}
func categoryHandler(router *mux.Router) {
	router.HandleFunc("/categories", CategoryController.GetCategory).Methods("GET")
	router.HandleFunc("/category/categoryId={categoryId}", CategoryController.GetSingleCategory).Methods("GET")
	router.HandleFunc("/addCategory", CategoryController.CreateCategory).Methods("POST", "OPTIONS")
	router.HandleFunc("/deleteCategory/categoryId={categoryId}", CategoryController.DeleteCategory).Methods("DELETE", "OPTIONS")
}
func cartHandler(router *mux.Router) {
	router.HandleFunc("/carts",CartController.GetCart).Methods("GET")
	router.Handle("/cart/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CartController.GetCartByCustomerId))).Methods("GET")
	router.Handle("/addToCart/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CartController.AddToCart))).Methods("POST", "OPTIONS")
	router.Handle("/removeAllItemsFromCart/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CartController.RemoveAllItemsFromCart))).Methods("POST", "OPTIONS")

	router.Handle("/cart/toggleProductToBuy/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CartController.ToggleBuyingProduct))).Methods("POST", "OPTIONS")
	router.Handle("/cart/increasingProductQuantity/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CartController.IncreasingProductQuantity))).Methods("POST", "OPTIONS")
	router.Handle("/cart/decreasingProductQuantity/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CartController.DecreasingProductQuantity))).Methods("POST", "OPTIONS")
}
func wishlistHandler(router *mux.Router) {
	router.HandleFunc("/wishlists",WishlistController.GetWishlist).Methods("GET")
	router.Handle("/wishlists/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(WishlistController.GetWishlistByCustomerId))).Methods("GET")
	router.Handle("/addToWishlist/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(WishlistController.AddToWishlist))).Methods("POST", "OPTIONS")
	router.Handle("/RemoveAllItemsFromWishlist/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(WishlistController.RemoveAllItemsFromWishlist))).Methods("POST", "OPTIONS")

	router.Handle("/wishlists/removeProductFromWishlist/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(WishlistController.RemoveProductFromWishlist))).Methods("POST", "OPTIONS")
}
func addressHandler(router *mux.Router) {
	router.HandleFunc("/addresses", AddressController.GetAddress).Methods("GET")
	router.Handle("/addAddress/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(AddressController.AddAddress))).Methods("POST", "OPTIONS")
	router.Handle("/updateAddress/addressId={addressId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(AddressController.UpdateAddress))).Methods("PUT", "OPTIONS")
	router.Handle("/deleteAddress/addressId={addressId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(AddressController.DeleteAddress))).Methods("DELETE", "OPTIONS")
}
func orderHandler(router *mux.Router) {
	router.HandleFunc("/orders", OrderController.GetOrder).Methods("GET")
	router.Handle("/orders/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.GetOrdersByCustomerId))).Methods("GET")
	router.Handle("/buyNow/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.BuyNow))).Methods("POST", "OPTIONS")
	router.Handle("/continueToPayment/orderId={orderId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.ContinueTOPayment))).Methods("POST", "OPTIONS")
	router.Handle("/replaceOrderRequest/orderId={orderId}&productId={productId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.ReplaceOrder))).Methods("PUT", "OPTIONS")
	router.Handle("/returnOrderRequest/orderId={orderId}&productId={productId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.ReturnOrder))).Methods("PUT", "OPTIONS")
	router.Handle("/cancelOrderRequest/orderId={orderId}&productId={productId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.CancelOrder))).Methods("PUT", "OPTIONS")

	router.Handle("/changeDeliveryDate/orderId={orderId}&productId={productId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(OrderController.ChangeDeliveryDate))).Methods("PUT", "OPTIONS")
}
func deliveryPartnerHandler(router *mux.Router) {
	router.HandleFunc("/deliverOrders", DeliveryPartnerController.GetDeliverOrders).Methods("GET")
	router.HandleFunc("/replaceOrders", DeliveryPartnerController.GetReplaceOrders).Methods("GET")
	router.HandleFunc("/returnOrders", DeliveryPartnerController.GetReturnOrders).Methods("GET")

	router.HandleFunc("/deliverOrderByDeliveryPartner/orderId={orderId}&productId={productId}&customerId={customerId}", DeliveryPartnerController.OrderDelivered).Methods("POST", "OPTIONS")
	router.HandleFunc("/replaceOrderByDeliveryPartner/orderId={orderId}&productId={productId}&customerId={customerId}", DeliveryPartnerController.OrderReplaced).Methods("POST", "OPTIONS")
	router.HandleFunc("/returnOrderByDeliveryPartner/orderId={orderId}&productId={productId}&customerId={customerId}", DeliveryPartnerController.OrderReturned).Methods("POST", "OPTIONS")
}
func reviewHandler(router *mux.Router) {
	router.HandleFunc("/reviews", ReviewController.GetReviews).Methods("GET")
	router.HandleFunc("/reviews/productId={productId}", ReviewController.GetReviewsByProductId).Methods("GET")
	router.HandleFunc("/mostRecentReviews/productId={productId}", ReviewController.GetReviewsOfProductByMostRecent).Methods("GET")
	router.HandleFunc("/topReviews/productId={productId}", ReviewController.GetReviewsOfProductByTopReviews).Methods("GET")

	router.Handle("/addReview/productId={productId}&customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ReviewController.AddReview))).Methods("POST", "OPTIONS")
}
func reportHandler(router *mux.Router) {
	router.Handle("/report/sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ReportController.ViewReports))).Methods("GET")
	router.Handle("/report/productId={productId}&sellerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ReportController.ViewReportByProductId))).Methods("GET")
}
func searchHistoryHandler(router *mux.Router) {
	router.HandleFunc("/searchHistory", SearchHistoryController.GetSearchHistory).Methods("GET")
	router.Handle("/searchHistory/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(SearchHistoryController.GetSingleSearchHistoryByCustomerId))).Methods("GET")
}
func customerCareHandler(router *mux.Router) {
	router.Handle("/customerCare/message/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(CustomerCareController.ChatWithMe))).Methods("POST", "OPTIONS")
}
func productRecommendationHandler(router *mux.Router) {
	router.Handle("/productRecommendation/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(ProductRecommendationController.GetRecommendatedProductByCustomerId))).Methods("GET")
}
func recentlyViewedProductHandler(router *mux.Router) {
	router.Handle("/recentlyViewedProduct/customerId={id}", auth.ProtectedMiddleware(http.HandlerFunc(RecentlyViewedProductController.GetRecentlyViewedProductByCustomerId))).Methods("GET")
}