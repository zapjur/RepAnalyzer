FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o db-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/db-service /app/db-service

EXPOSE 50051
CMD ["/app/db-service"]
