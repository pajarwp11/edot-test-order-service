package order

import (
	"order-service/models/cart"
	"time"
)

type Order struct {
	Id         int       `db:"id" json:"id"`
	UserId     int       `db:"user_id" json:"user_id"`
	TotalPrice int       `db:"total_price" json:"total_price"`
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
