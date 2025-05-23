package service

import (
	"context"
	"database/sql"
	"github.com/mirawbtw/assignment6/internal/domain"
	"github.com/mirawbtw/assignment6/internal/repository"
	"time"
)

type ProductService interface {
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProductStock(ctx context.Context, id string, quantity int) error
}

type productService struct {
	productRepo repository.ProductRepository
	redisRepo   repository.RedisRepository
	db          *sql.DB
}

func NewProductService(
	productRepo repository.ProductRepository,
	redisRepo repository.RedisRepository,
	db *sql.DB,
) ProductService {
	return &productService{
		productRepo: productRepo,
		redisRepo:   redisRepo,
		db:          db,
	}
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	if product, err := s.redisRepo.GetProduct(ctx, id); err == nil {
		return product, nil
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.redisRepo.SetProduct(ctx, product, 10*time.Minute)

	return product, nil
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}

func (s *productService) UpdateProductStock(ctx context.Context, id string, quantity int) error {
	// Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txRepo := s.productRepo.WithTx(tx)

	if err := txRepo.UpdateStock(ctx, id, quantity); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	_ = s.redisRepo.DeleteProduct(ctx, id)

	return nil
}
