package main

import (
	"fmt"
	"log"
	"net/http"
	"toko-online-go/config"
	"toko-online-go/controllers"
	"toko-online-go/middlewares"

	"github.com/rs/cors"
)

func main() {
	config.ConnectDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/products", middlewares.Auth(controllers.ProductsHandler))
	mux.HandleFunc("/register", controllers.Register)
	mux.HandleFunc("/login", controllers.Login)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	fmt.Println("Server berjalan di http://localhost:8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
