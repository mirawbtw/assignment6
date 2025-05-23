package repository

import (
	"context"
	"database/sql"
	"github.com/mirawbtw/assignment6/internal/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	UpdateStock(ctx context.Context, id string, quantity int) error
	WithTx(tx *sql.Tx) ProductRepository
}

type productRepository struct {
	db *sql.DB
}

type productTxRepository struct {
	tx *sql.Tx
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, description, price, stock FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO products (id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5)",
		product.ID, product.Name, product.Description, product.Price, product.Stock)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE products SET stock = $1 WHERE id = $2",
		quantity, id)
	return err
}

func (r *productRepository) WithTx(tx *sql.Tx) ProductRepository {
	return &productTxRepository{tx: tx}
}

func (r *productTxRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product := &domain.Product{}
	err := r.tx.QueryRowContext(ctx,
		"SELECT id, name, description, price, stock FROM products WHERE id = $1", id).
		Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productTxRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.tx.ExecContext(ctx,
		"INSERT INTO products (id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5)",
		product.ID, product.Name, product.Description, product.Price, product.Stock)
	return err
}

func (r *productTxRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	_, err := r.tx.ExecContext(ctx,
		"UPDATE products SET stock = $1 WHERE id = $2",
		quantity, id)
	return err
}

func (r *productTxRepository) WithTx(tx *sql.Tx) ProductRepository {
	return &productTxRepository{tx: tx}
}
