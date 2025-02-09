package order

import (
	"encoding/json"
	"errors"
	"order-service/models/order"
)

func (o *OrderHandler) UpdateStatusByConsumer(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := order.UpdateStatusRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = o.orderUsecase.UpdateStatus(&request, true)
	if err != nil {
		return err
	}
	return nil
}
