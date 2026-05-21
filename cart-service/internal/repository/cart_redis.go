package repository

import (
	"context"
	"encoding/json"
	"time"

	"cart/internal/cart"
	"github.com/go-redis/redis/v8"
)

type CartRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewCartRedisRepository(client *redis.Client) *CartRedisRepository {
	return &CartRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *CartRedisRepository) Save(cart *cart.Cart) error {
	data, _ := json.Marshal(cart)
	return r.client.Set(r.ctx, "cart:"+cart.UserID, data, 30*time.Minute).Err()
}

func (r *CartRedisRepository) FindByID(userID string) (*cart.Cart, error) {
	data, err := r.client.Get(r.ctx, "cart:"+userID).Result()
	if err != nil {
		return &cart.Cart{UserID: userID, Items: []cart.Item{}}, nil
	}

	var cart cart.Cart
	_ = json.Unmarshal([]byte(data), &cart)
	return &cart, nil
}