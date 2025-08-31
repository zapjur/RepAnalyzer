FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY ../.. .

RUN go mod tidy
RUN go build -o analysis-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/analysis-service /app/analysis-service

CMD ["/app/analysis-service"]
