# Product Service 

Микросервис для управления продуктами с gRPC API, Redis-кэшированием и PostgreSQL

## Функционал
- gRPC API:
  - Создание продукта (`CreateProduct`)
  - Получение продукта с кэшированием (`GetProduct`)
  - Обновление остатков с транзакцией (`UpdateProductStock`)
  
- Технологии:
  - Go 1.20+
  - gRPC
  - PostgreSQL
  - Redis
  - Protocol Buffers

## Запуск проекта

### Требования
- Установленные:
  - Go 1.20+
  - PostgreSQL 14+
  - Redis 6+
  - protoc (компилятор Protocol Buffers)

### 1. Клонирование репозитория
```bash
git clone https://github.com/mirawbtw/assignment6-product-service.git
cd assignment6-product-service
