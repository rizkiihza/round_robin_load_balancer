# Golang round robin

Golang round robin is a load balancer implemented using golang

## Why golang
Golang has very robust standard library, this project is made without using external dependency. That is except of gomock, but we only use gomock for unit testing

## How to start

```bash
docker-compose up --build
```

## Usage

```bash
curl --location 'localhost:9000' \
--header 'Content-Type: application/json' \
--data '{
    "Hello": "world"
}'
```