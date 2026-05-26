FROM golang:1.23-bookworm AS build
WORKDIR /app/api-gateway
COPY api-gateway/ ./
COPY pkg ../pkg
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o /app/api-gateway .

FROM debian:bookworm-slim
COPY --from=build /app/api-gateway /usr/local/bin/api-gateway
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api-gateway"]