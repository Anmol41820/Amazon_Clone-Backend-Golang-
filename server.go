package main

import(
	router "Amazon_Server/Router"
	"fmt"
	"net/http"
)

func main(){
	router := router.GetRouter()
	fmt.Println("Server is running : 3000")
	http.ListenAndServe(":3000",router)
}