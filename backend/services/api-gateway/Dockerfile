FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o api-gateway ./cmd

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y ffmpeg ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /app/api-gateway /app/api-gateway

EXPOSE 8080
CMD ["/app/api-gateway"]
