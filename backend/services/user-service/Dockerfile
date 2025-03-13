FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o user-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/user-service /app/user-service

EXPOSE 50051
CMD ["/app/user-service"]
