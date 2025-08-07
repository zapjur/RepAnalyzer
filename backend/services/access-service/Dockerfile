FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o access-service ./cmd

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /app/access-service /app/access-service

CMD ["/app/access-service"]
