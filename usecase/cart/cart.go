package cart

import (
	"order-service/models/cart"
)

type CartRepository interface {
	Insert(cart *cart.Cart) error
	Get(userId int) (*cart.Cart, error)
}

type CartUsecase struct {
	cartRepo CartRepository
}

func NewCartUsecase(cartRepo CartRepository) *CartUsecase {
	return &CartUsecase{
		cartRepo: cartRepo,
	}
}

func (c *CartUsecase) Insert(cart *cart.Cart) error {
	return c.cartRepo.Insert(cart)
}

func (c *CartUsecase) Get(userId int) (*cart.Cart, error) {
	return c.cartRepo.Get(userId)
}
