package order

import (
	"encoding/json"
	"net/http"
	"order-service/models/order"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type OrderUsecase interface {
	Checkout(orderCheckout *order.CheckoutRequest) error
	UpdateStatus(updateStatus *order.UpdateStatusRequest, isConsumer bool) error
}

type OrderHandler struct {
	orderUsecase OrderUsecase
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var validate = validator.New()

func NewOrderHandler(orderUsecase OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (o *OrderHandler) Checkout(w http.ResponseWriter, req *http.Request) {
	request := order.CheckoutRequest{}
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

	userId := req.Header.Get("X-User-ID")
	request.UserId, _ = strconv.Atoi(userId)

	err := o.orderUsecase.Checkout(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "order created"
	json.NewEncoder(w).Encode(response)
}

func (o *OrderHandler) UpdateStatus(w http.ResponseWriter, req *http.Request) {
	request := order.UpdateStatusRequest{}
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

	vars := mux.Vars(req)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "id is required"
		json.NewEncoder(w).Encode(response)
		return
	}
	var err error
	request.Id, err = strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "id must be numeric"
		json.NewEncoder(w).Encode(response)
		return
	}

	err = o.orderUsecase.UpdateStatus(&request, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "order status updated"
	json.NewEncoder(w).Encode(response)
}
