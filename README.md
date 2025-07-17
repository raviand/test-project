# Go REST API

This is a REST API built with Go that manages products and users. It uses Chi for routing, MySQL for product data, and DynamoDB for user data.

## Features

-   CRUD operations for products and users
-   Token-based authentication
-   Audit logging of requests

## Prerequisites

-   Go
-   Docker
-   MySQL
-   DynamoDB Local

## Getting Started

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/raviand/test-project.git
    ```

2.  **Set up the environment variables:**

    Create a `.env` file in the root of the project and add the following variables:

    ```
    MYSQL_USER=your_mysql_user
    MYSQL_PASSWORD=your_mysql_password
    DYNAMO_ENDPOINT=http://localhost:8000
    TOKEN=your_secret_token
    ```

3.  **Run the database:**

    ```bash
    docker-compose up -d
    ```

4.  **Run the application:**

    ```bash
    go run cmd/main.go
    ```

## API Endpoints

### Products

#### Get all products

```bash
curl -X GET http://localhost:8080/product -H "Authorization: your_secret_token"
```

#### Get a product by ID

```bash
curl -X GET http://localhost:8080/product/1 -H "Authorization: your_secret_token"
```

#### Create a new product

```bash
curl -X POST http://localhost:8080/product -H "Authorization: your_secret_token" -d '{
    "name": "Macbook Pro",
    "price": 1500.50,
    "quantity": 10,
    "code_value": "mbp-2023",
    "is_published": true,
    "expiration": "25/12/2025",
    "product_type": 1
}'
```

#### Update or create a product

```bash
curl -X PUT http://localhost:8080/product/1 -H "Authorization: your_secret_token" -d '{
    "name": "Macbook Pro M2",
    "price": 1600.50,
    "quantity": 10,
    "code_value": "mbp-2023-m2",
    "is_published": true,
    "expiration": "25/12/2025",
    "product_type": 1
}'
```

#### Partially update a product

```bash
curl -X PATCH http://localhost:8080/product/1 -H "Authorization: your_secret_token" -d '{
    "price": 1650.50
}'
```

#### Delete a product

```bash
curl -X DELETE http://localhost:8080/product/1 -H "Authorization: your_secret_token"
```

### Users

#### Get all users

```bash
curl -X GET http://localhost:8080/user -H "Authorization: your_secret_token"
```

#### Get a user by ID

```bash
curl -X GET http://localhost:8080/user/1 -H "Authorization: your_secret_token"
```

#### Create a new user

```bash
curl -X POST http://localhost:8080/user -H "Authorization: your_secret_token" -d '{
    "id": "1",
    "firstname": "John",
    "lastname": "Doe",
    "username": "johndoe",
    "password": "password",
    "email": "johndoe@example.com",
    "ip": "127.0.0.1",
    "macAddress": "00:00:00:00:00:00",
    "website": "johndoe.com",
    "image": "johndoe.png"
}'
```

#### Delete a user

```bash
curl -X DELETE http://localhost:8080/user/1 -H "Authorization: your_secret_token"
```
