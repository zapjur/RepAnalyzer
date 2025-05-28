FROM golang:1.24.1-alpine AS build

WORKDIR /app
COPY ../.. .

RUN go mod tidy
RUN go build -o orchestrator ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/orchestrator /app/orchestrator

EXPOSE 50051
CMD ["/app/orchestrator"]
