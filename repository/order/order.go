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

func (o *OrderRepository) Insert(tx *sqlx.Tx, order *order.Order) (int, error) {
	res, err := tx.Exec("INSERT INTO orders (user_id,total_price,status) VALUES (?,?,?)", order.UserId, order.TotalPrice, order.Status)
	if err != nil {
		return 0, err
	}
	orderId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(orderId), nil
}

func (o *OrderRepository) InsertDetails(tx *sqlx.Tx, orderId int, products []cart.Product) error {
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

	_, err := tx.Exec(query, args...)
	return err
}

func (o *OrderRepository) UpdateStatus(tx *sqlx.Tx, id int, status string) error {
	_, err := o.mysql.Exec("UPDATE orders SET status=? WHERE id=?", status, id)
	return err
}

func (o *OrderRepository) GetOrderWithDetails(orderId int) (*order.OrderWithDetails, error) {
	query := `
	SELECT 
		o.id, o.user_id, o.total_price, o.status,
		od.id as order_detail_id, od.product_id, od.product_name, od.quantity, od.price, od.total_price
	FROM orders o
	LEFT JOIN order_details od ON o.id = od.order_id
	WHERE o.id = ?
	`

	rows, err := o.mysql.Queryx(query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderWithDetail order.OrderWithDetails
	orderWithDetail.Details = []order.OrderDetail{}

	for rows.Next() {
		var detail order.OrderDetail
		err := rows.Scan(
			&orderWithDetail.ID, &orderWithDetail.UserId, &orderWithDetail.TotalPrice, &orderWithDetail.Status,
			&detail.ID, &detail.ProductId, &detail.ProductName, &detail.Quantity, &detail.Price, &detail.TotalPrice,
		)
		if err != nil {
			return nil, err
		}

		if detail.ID != 0 {
			orderWithDetail.Details = append(orderWithDetail.Details, detail)
		}
	}

	if orderWithDetail.ID == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return &orderWithDetail, nil
}
