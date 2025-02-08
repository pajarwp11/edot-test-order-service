package cart

import (
	"encoding/json"
	"net/http"
	"order-service/models/cart"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type CartUsecase interface {
	Insert(cart *cart.Cart) error
	Get(userId int) (*cart.Cart, error)
}

type CartHandler struct {
	cartUsecase CartUsecase
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var validate = validator.New()

func NewCartHandler(cartUsecase CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase: cartUsecase,
	}
}

func (c *CartHandler) Insert(w http.ResponseWriter, req *http.Request) {
	request := cart.Cart{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := c.cartUsecase.Insert(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "cart inserted"
	json.NewEncoder(w).Encode(response)
}

func (c *CartHandler) Get(w http.ResponseWriter, req *http.Request) {
	response := Response{}
	w.Header().Set("Content-Type", "application/json")

	userId := req.Header.Get("X-User-ID")
	userIdInt, _ := strconv.Atoi(userId)

	cart, err := c.cartUsecase.Get(userIdInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Message = "get cart success"
	response.Data = cart
	json.NewEncoder(w).Encode(response)
}
