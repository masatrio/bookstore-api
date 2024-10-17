# Online Bookstore API

An API for an online bookstore allowing customers to create accounts, view a list of books, make orders, and view their order history.

---

## **Table of Contents**

- [Features](#features)
- [Technologies](#technologies)
- [API Docs](#api-docs)
- [Project Structure](#project-structure)
- [Database Schema](#database-schemas)
- [Setup and Installation](#setup-and-installation)
- [Running Tests](#running-tests)

---

## **Features**

- **Create Customer Account**: Sign up for an account using a unique email.
- **View Books**: Browse the available books.
- **Place Orders**: Make an order with multiple books.
- **View Order History**: See all previous orders.

---

## **Technologies**

- **Programming Language**: `Golang`
- **Database**: `PostgreSQL`
- **Cache**: `Redis` ( not yet implemented )
- **Search**: `Elasticsearch` ( not yet implemented )
- **Observability Framework**: `Open Telemetry`

---

## **API Docs**

[https://app.swaggerhub.com/apis/masatrio/bookstore-api/1.0.0](#https://app.swaggerhub.com/apis/masatrio/bookstore-api/1.0.0)

---

## **Project Structure**

```/bookstore-api
│
├── Dockerfile
├── README.md
├── docker-compose.yml
├── go.mod
├── go.sum
├── sample.env
├── otel-collector-config.yaml
│
├── /cmd
│   ├── /migrate
│   │   └── main.go  # database migrations
│   ├── /seed
│   │   └── main.go  # data seeding
│   └── /server
│       └── main.go  # main server entry point
│
├── /config
│   └── config.go  # application configuration
│
├── /internal
│   ├── /delivery
│   │   └── /http
│   │       ├── handlers.go  # HTTP request handlers
│   │       ├── routes.go  # route definitions
│   │       └── /middleware
│   │           ├── jwt.go  # JWT authentication middleware
│   │           ├── otel.go  # OpenTelemetry integration
│   │           └── panic.go  # panic recovery middleware
│   │
│   ├── /domain
│   │   ├── /cache
│   │   │   ├── book_cache.go  # book caching interface
│   │   │   ├── customer_cache.go  # customer caching interface
│   │   │   └── order_cache.go  # order caching interface
│   │   ├── /delivery
│   │   │   └── http.go  # delivery interface
│   │   ├── /repository
│   │   │   ├── book_repository.go  # book repository interface
│   │   │   ├── order_repository.go  # order repository interface
│   │   │   ├── repository.go  # common repository interface
│   │   │   └── user_repository.go  # user repository interface
│   │   └── /usecase
│   │       ├── book_usecase.go  # book use case logic
│   │       ├── order_usecase.go  # order use case logic
│   │       └── user_usecase.go  # user use case logic
│   │
│   ├── /repository
│   │   ├── /cache
│   │   │   └── /redis
│   │   │       ├── book_cache.go  # Redis book cache implementation
│   │   │       ├── customer_cache.go  # Redis customer cache implementation
│   │   │       └── order_cache.go  # Redis order cache implementation
│   │   ├── /db
│   │   │   └── /postgresql
│   │   │       ├── book_repository.go  # PostgreSQL book repository
│   │   │       ├── order_item_repository.go  # PostgreSQL order item repository
│   │   │       ├── order_repository.go  # PostgreSQL order repository
│   │   │       ├── postgresql.go  # common PostgreSQL setup
│   │   │       ├── repository.go  # common repository implementation
│   │   │       └── user_repository.go  # PostgreSQL user repository
│   │   └── /search
│   │       └── /elasticsearch
│   │           └── search.go  # Elasticsearch search implementation
│   │
│   └── /usecase
│       ├── /book
│       │   └── book.go  # book use case implementation
│       ├── /order
│       │   └── order.go  # order use case implementation
│       └── /user
│           └── user.go  # user use case implementation
│
├── /migrations  # SQL migration files
│   ├── 1_create_users_table.up.sql
│   ├── 1_create_users_table.down.sql
│   ├── 2_create_books_table.up.sql
│   ├── 2_create_books_table.down.sql
│   ├── 3_create_orders_table.up.sql
│   ├── 3_create_orders_table.down.sql
│   ├── 4_create_order_items_table.up.sql
│   └── 4_create_order_items_table.down.sql
│
└── /utils
    ├── db.go  # database utility functions
    ├── errors.go  # error handling utilities
    ├── jwt.go  # JWT utility functions
    └── tracer.go  # tracing utility functions
```
---

## **Database Schemas**
- **Users Table**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
- **Books Table**
```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
- **Orders Table**
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
- **OrderItems Table**
```sql
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    book_id INT NOT NULL,
    quantity INT NOT NULL
);
```
---

## **Setup and Installation**

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/masatrio/bookstore-api.git
    ```

2. **Prepare ENV**: Prepare the env by creating .env file, or copy the sample.env
    ```bash
    mv sample.env .env
    ```

3. **Setup Docker**: Ensure Docker and Docker Compose are installed. Then run:
    ```bash
    docker-compose up --build
    ```

3. **Access API Documentation**: 
   Visit [https://app.swaggerhub.com/apis/masatrio/bookstore-api/1.0.0](#https://app.swaggerhub.com/apis/masatrio/bookstore-api/1.0.0) to access Swagger API documentation.

---

## **Running Tests**

Run the tests with the following command:
```bash
mockgen -source=./internal/domain/usecase/user_usecase.go -destination=./internal/domain/usecase/mocks/user_usecase_mock.go -package=mocks
mockgen -source=./internal/domain/usecase/book_usecase.go -destination=./internal/domain/usecase/mocks/book_usecase_mock.go -package=mocks
mockgen -source=./internal/domain/usecase/order_usecase.go -destination=./internal/domain/usecase/mocks/order_usecase_mock.go -package=mocks
go test ./...
```
---