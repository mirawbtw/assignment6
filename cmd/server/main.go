package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	"github.com/mirawbtw/assignment6/api"
	"github.com/mirawbtw/assignment6/internal/config"
	"github.com/mirawbtw/assignment6/internal/handler"
	"github.com/mirawbtw/assignment6/internal/repository"
	"github.com/mirawbtw/assignment6/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Postgres.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	productRepo := repository.NewProductRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)

	productService := service.NewProductService(productRepo, redisRepo, db)

	server := grpc.NewServer(
		grpc.ConnectionTimeout(cfg.GRPC.Timeout),
	)
	api.RegisterProductServiceServer(server, handler.NewProductHandler(productService))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server started on :%d", cfg.GRPC.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
