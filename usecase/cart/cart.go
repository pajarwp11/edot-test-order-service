package cart

import (
	"errors"
	"order-service/models/cart"

	"github.com/redis/go-redis/v9"
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
	cartData, err := c.cartRepo.Get(userId)
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("user does not have cart")
		}
		return nil, err
	}
	return cartData, nil
}
