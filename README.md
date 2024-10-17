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
├── /cmd
│   └── /server
│       └── main.go #main server
│
├── /internal
│   ├── /domain
│   │   ├── /usecase #usecase interface
│   │   ├── /repository #repository interface
│   │   ├── /delivery #delivery interface
│   │   └── /cache #cache interface
│   │
│   ├── /usecase #usecase implementation
│   │   ├── /book
│   │   ├── /order
│   │   └── /user
│   │
│   ├── /repository
│   │   ├── /db #database implementation
│   │   │   └── /postgresql
│   │   └── /cache #cache implementation
│   │       └── /redis
│   │
│   └── /delivery #delivery implementation
│       └── /http
│
├── /middleware #middleware
│   ├── jwt.go
│   └── tracing.go
│
└── /config
    └── config.go
```
---

## **Database Schemas**
- **Customers Table**
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

2. **Setup Docker**: Ensure Docker and Docker Compose are installed. Then run:
    ```bash
    docker-compose up --build
    ```

3. **Run Migrations**: The project includes a database migration tool for PostgreSQL:
    ```bash
    go run cmd/migrate/main.go
    ```

4. **Run the Application**:
    ```bash
    go run cmd/server/main.go
    ```

5. **Access API Documentation**: 
   Visit [https://app.swaggerhub.com/apis/masatrio/bookstore-api/1.0.0](#https://app.swaggerhub.com/apis/masatrio/bookstore-api/1.0.0) to access Swagger API documentation.

---

## **Running Tests**

Run the tests with the following command:
```bash
go test ./...
```
---