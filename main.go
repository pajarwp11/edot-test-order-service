package main

import (
	"fmt"
	"log"
	"net/http"
	"order-service/db/mysql"
	"order-service/db/redis"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	redis.Connect()
	router := mux.NewRouter()

	fmt.Println("server is running")
	err := http.ListenAndServe(":8004", router)
	if err != nil {
		log.Fatal(err)
	}
}
