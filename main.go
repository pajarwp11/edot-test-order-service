package main

import (
	"fmt"
	"log"
	"net/http"
	"order-service/db/mysql"
	"order-service/db/redis"

	cartHandler "order-service/handler/cart"
	"order-service/middleware"
	cartRepo "order-service/repository/cart"
	cartUsecase "order-service/usecase/cart"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	redis.Connect()
	router := mux.NewRouter()
	cartRepository := cartRepo.NewCartRepository(redis.Redis)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepository)
	cartHandler := cartHandler.NewCartHandler(cartUsecase)
	router.Handle("/cart/insert", middleware.JWTMiddleware(http.HandlerFunc(cartHandler.Insert))).Methods(http.MethodPost)
	router.Handle("/cart", middleware.JWTMiddleware(http.HandlerFunc(cartHandler.Get))).Methods(http.MethodGet)

	fmt.Println("server is running")
	err := http.ListenAndServe(":8004", router)
	if err != nil {
		log.Fatal(err)
	}
}
