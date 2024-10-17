
FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main ./cmd/server

FROM alpine:3.18
WORKDIR /app

COPY --from=build /app/main .
COPY .env .env

EXPOSE 8080
CMD ["./main"]
