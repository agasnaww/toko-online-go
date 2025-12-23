package main

import (
	"fmt"
	"log"
	"net/http"

	"toko-online-go/config"
	"toko-online-go/controllers"
)

func main() {
	config.ConnectDB()
	http.HandleFunc("/products", controllers.ProductsHandler)
	fmt.Println("Server berjalan di http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
