package order

import (
	"order-service/models/cart"
	"time"
)

type Order struct {
	Id         int       `db:"id" json:"id"`
	UserId     int       `db:"user_id" json:"user_id"`
	TotalPrice int       `db:"total_price" json:"total_price"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type CheckoutRequest struct {
	UserId   int            `json:"user_id"`
	IsCart   *int           `json:"is_cart" validate:"required"`
	Products []cart.Product `json:"products" validate:"required"`
}

type StockOperationRequest struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type UpdateStatusRequest struct {
	Id     int    `json:"status"`
	Status string `json:"status" validate:"required"`
}

type OrderWithDetails struct {
	ID         int           `db:"id" json:"id"`
	UserId     int           `db:"user_id" json:"user_id"`
	TotalPrice int           `db:"total_price" json:"total_price"`
	Status     string        `db:"status" json:"status"`
	Details    []OrderDetail `json:"details"`
}

type OrderDetail struct {
	ID          int    `db:"id" json:"id"`
	OrderID     int    `db:"order_id" json:"order_id"`
	ProductId   int    `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Quantity    int    `db:"quantity" json:"quantity"`
	Price       int    `db:"price" json:"price"`
	TotalPrice  int    `db:"total_price" json:"total_price"`
}
