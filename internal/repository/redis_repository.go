package repository

import (
	"context"
	"encoding/json"
	"github.com/mirawbtw/assignment6/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository interface {
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	SetProduct(ctx context.Context, product *domain.Product, expiration time.Duration) error
	DeleteProduct(ctx context.Context, id string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	data, err := r.client.Get(ctx, "product:"+id).Bytes()
	if err != nil {
		return nil, err
	}

	var product domain.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *redisRepository) SetProduct(ctx context.Context, product *domain.Product, expiration time.Duration) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "product:"+product.ID, data, expiration).Err()
}

func (r *redisRepository) DeleteProduct(ctx context.Context, id string) error {
	return r.client.Del(ctx, "product:"+id).Err()
}
