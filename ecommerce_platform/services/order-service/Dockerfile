FROM golang:1.23-bookworm AS build
WORKDIR /app/order-service
COPY order-service/ ./
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o /app/order-service .

FROM debian:bookworm-slim
COPY --from=build /app/order-service /usr/local/bin/order-service
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/order-service"]
