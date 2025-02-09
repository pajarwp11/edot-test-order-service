package order

import (
	"fmt"
	"order-service/models/cart"
	"order-service/models/order"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	mysql *sqlx.DB
}

func NewOrderRepository(mysql *sqlx.DB) *OrderRepository {
	return &OrderRepository{
		mysql: mysql,
	}
}

func (p *OrderRepository) Insert(order *order.Order) error {
	_, err := p.mysql.Exec("INSERT INTO orders (user_id,total_price) VALUES (?,?)", order.UserId, order.TotalPrice)
	return err
}

func (p *OrderRepository) InsertDetails(orderId int, products []*cart.Product) error {
	if len(products) == 0 {
		return nil
	}

	values := make([]string, 0, len(products))
	args := make([]interface{}, 0, len(products)*6)

	for _, product := range products {
		values = append(values, "(?,?,?,?,?,?)")
		args = append(args, orderId, product.ProductId, product.ProductName, product.Quantity, product.Price, product.TotalPrice)
	}

	query := fmt.Sprintf(
		"INSERT INTO order_details (order_id, product_id, product_name, quantity, price, total_price) VALUES %s",
		strings.Join(values, ","),
	)

	_, err := p.mysql.Exec(query, args...)
	return err
}
