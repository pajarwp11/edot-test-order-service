package main

import (
	"fmt"
	"log"
	"net/http"
	"order-service/conn/mysql"
	"order-service/conn/rabbitmq"
	"order-service/conn/redis"

	cartHandler "order-service/handler/cart"
	"order-service/middleware"
	cartRepo "order-service/repository/cart"
	cartUsecase "order-service/usecase/cart"

	orderHandler "order-service/handler/order"
	orderRepo "order-service/repository/order"
	orderUsecase "order-service/usecase/order"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	redis.Connect()
	rabbitmq.Connect()
	rabbitPublisher := rabbitmq.NewRabbitPublisher(rabbitmq.RabbitConn)

	router := mux.NewRouter()
	cartRepository := cartRepo.NewCartRepository(redis.Redis)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepository)
	cartHandler := cartHandler.NewCartHandler(cartUsecase)
	router.Handle("/cart/insert", middleware.JWTMiddleware(http.HandlerFunc(cartHandler.Insert))).Methods(http.MethodPost)
	router.Handle("/cart", middleware.JWTMiddleware(http.HandlerFunc(cartHandler.Get))).Methods(http.MethodGet)

	orderRepository := orderRepo.NewOrderRepository(mysql.MySQL)
	orderUsecase := orderUsecase.NewOrderUsecase(orderRepository, cartRepository, mysql.MySQL, rabbitPublisher)
	orderHandler := orderHandler.NewOrderHandler(orderUsecase)
	router.Handle("/order/checkout", middleware.JWTMiddleware(http.HandlerFunc(orderHandler.Checkout))).Methods(http.MethodPost)
	router.Handle("/order/update-status", middleware.JWTMiddleware(http.HandlerFunc(orderHandler.UpdateStatus))).Methods(http.MethodPut)

	rabbitConsumer := rabbitmq.NewRabbitConsumer(rabbitmq.RabbitConn, orderHandler)
	go rabbitConsumer.ConsumeEvents()

	fmt.Println("server is running")
	err := http.ListenAndServe(":8004", router)
	if err != nil {
		log.Fatal(err)
	}
}
