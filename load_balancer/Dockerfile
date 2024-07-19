# build container
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# app container
FROM alpine:latest
ARG APP_PORT
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE ${APP_PORT}
CMD ["./main"]
