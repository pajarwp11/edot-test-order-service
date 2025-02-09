package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/models/cart"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartRepository struct {
	redisClient *redis.Client
}

func NewCartRepository(redis *redis.Client) *CartRepository {
	return &CartRepository{
		redisClient: redis,
	}
}

func (c *CartRepository) Insert(cart *cart.Cart) error {
	cartByte, err := json.Marshal(cart)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("cart:%d", cart.UserId)
	ctx := context.Background()
	err = c.redisClient.Set(ctx, key, string(cartByte), 30*24*time.Hour).Err()
	return err
}

func (c *CartRepository) Get(userId int) (*cart.Cart, error) {
	key := fmt.Sprintf("cart:%d", userId)
	ctx := context.Background()
	cartData, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	cart := cart.Cart{}
	err = json.Unmarshal([]byte(cartData), &cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (c *CartRepository) Delete(userId int) error {
	key := fmt.Sprintf("cart:%d", userId)
	ctx := context.Background()
	return c.redisClient.Del(ctx, key).Err()
}
