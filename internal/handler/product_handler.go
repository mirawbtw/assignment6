package handler

import (
	"context"
	"log"

	"github.com/mirawbtw/assignment6/api"
	"github.com/mirawbtw/assignment6/internal/domain"
	"github.com/mirawbtw/assignment6/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	api.UnimplementedProductServiceServer
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		UnimplementedProductServiceServer: api.UnimplementedProductServiceServer{},
		service:                           service,
	}
}
func (h *ProductHandler) GetProduct(ctx context.Context, req *api.GetProductRequest) (*api.ProductResponse, error) {
	product, err := h.service.GetProduct(ctx, req.Id)
	if err != nil {
		log.Printf("GetProduct error: %v", err)
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &api.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *api.CreateProductRequest) (*api.ProductResponse, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := h.service.CreateProduct(ctx, product); err != nil {
		log.Printf("CreateProduct error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create product")
	}

	return &api.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (h *ProductHandler) UpdateProductStock(ctx context.Context, req *api.UpdateStockRequest) (*api.ProductResponse, error) {
	if err := h.service.UpdateProductStock(ctx, req.Id, int(req.Quantity)); err != nil {
		log.Printf("UpdateProductStock error: %v", err)
		return nil, status.Error(codes.Internal, "failed to update product stock")
	}

	product, err := h.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found after update")
	}

	return &api.ProductResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}
