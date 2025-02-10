package order

import (
	"fmt"
	"order-service/models/cart"
	"order-service/models/order"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Insert(tx *sqlx.Tx, order *order.Order) (int, error)
	InsertDetails(tx *sqlx.Tx, orderId int, products []cart.Product) error
	UpdateStatus(tx *sqlx.Tx, id int, status string) error
	GetOrderWithDetails(orderId int) (*order.OrderWithDetails, error)
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

	reserveStocks := order.StockOperationOrderRequest{}
	totalPrice := 0
	for _, product := range orderCheckout.Products {
		totalPrice += product.TotalPrice
		reserveStock := order.StockOperationRequest{
			ProductId: product.ProductId,
			Quantity:  product.Quantity,
		}
		reserveStocks.StockOperations = append(reserveStocks.StockOperations, reserveStock)
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

	reserveStocks.OrderId = orderId
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

func (o *OrderUsecase) UpdateStatus(updateStatus *order.UpdateStatusRequest) error {
	orderWithDetail, err := o.orderRepo.GetOrderWithDetails(updateStatus.Id)
	if err != nil {
		return err
	}
	if orderWithDetail.Status != "pending" {
		return fmt.Errorf("status is already %s", orderWithDetail.Status)
	}
	tx, err := o.mysql.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = o.orderRepo.UpdateStatus(tx, updateStatus.Id, updateStatus.Status)
	if err != nil {
		return err
	}

	var event string
	if updateStatus.Status == "success" {
		event = "stock.release"
	} else {
		event = "stock.return"
	}

	stockOperation := order.StockOperationOrderRequest{}
	stockOperation.OrderId = updateStatus.Id
	for _, product := range orderWithDetail.Details {
		reserveStock := order.StockOperationRequest{
			ProductId: product.ProductId,
			Quantity:  product.Quantity,
		}
		stockOperation.StockOperations = append(stockOperation.StockOperations, reserveStock)
	}

	err = o.publisher.PublishEvent(event, stockOperation)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderUsecase) UpdateStatusConsumer(updateStatus *order.UpdateStatusRequest) error {
	orderWithDetail, err := o.orderRepo.GetOrderWithDetails(updateStatus.Id)
	if err != nil {
		return err
	}
	if orderWithDetail.Status != "pending" {
		return fmt.Errorf("status is already %s", orderWithDetail.Status)
	}
	tx, err := o.mysql.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = o.orderRepo.UpdateStatus(tx, updateStatus.Id, updateStatus.Status)
	if err != nil {
		return err
	}

	return nil
}
