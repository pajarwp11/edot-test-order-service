package order

import (
	"order-service/models/cart"
	"order-service/models/order"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Insert(tx *sqlx.Tx, order *order.Order) (int, error)
	InsertDetails(tx *sqlx.Tx, orderId int, products []cart.Product) error
}

type CartRepository interface {
	Delete(userId int) error
}

type Publisher interface {
	PublishEvent(eventType string, data interface{}) error
}

type OrderUsecase struct {
	orderRepo      OrderRepository
	cartRepository CartRepository
	mysql          *sqlx.DB
	publisher      Publisher
}

func NewOrderUsecase(orderRepo OrderRepository, cartRepository CartRepository, mysql *sqlx.DB, publisher Publisher) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:      orderRepo,
		cartRepository: cartRepository,
		mysql:          mysql,
		publisher:      publisher,
	}
}

func (o *OrderUsecase) Checkout(orderCheckout *order.CheckoutRequest) error {
	tx, err := o.mysql.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	reserveStocks := []order.StockOperationRequest{}
	totalPrice := 0
	for _, product := range orderCheckout.Products {
		totalPrice += product.TotalPrice
		reserveStock := order.StockOperationRequest{
			ProductId: product.ProductId,
			Quantity:  product.Quantity,
		}
		reserveStocks = append(reserveStocks, reserveStock)
	}

	orderData := order.Order{
		UserId:     orderCheckout.UserId,
		TotalPrice: totalPrice,
		Status:     "pending",
	}
	orderId, err := o.orderRepo.Insert(tx, &orderData)
	if err != nil {
		return err
	}
	err = o.orderRepo.InsertDetails(tx, orderId, orderCheckout.Products)
	if err != nil {
		return err
	}

	err = o.publisher.PublishEvent("stock.reserve", reserveStocks)
	if err != nil {
		return err
	}
	tx.Commit()
	if orderCheckout.IsCart != nil && *orderCheckout.IsCart == 1 {
		o.cartRepository.Delete(orderCheckout.UserId)
	}
	return nil
}
