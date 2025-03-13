FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o api-gateway ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/api-gateway /app/api-gateway

EXPOSE 8080
CMD ["/app/api-gateway"]
