# FROM golang:1.23.4-alpine AS builder

# WORKDIR /app
# COPY . .

# RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
# COPY --from=builder /app/main .
COPY main /app/main
# COPY --from=builder /app/.env .

EXPOSE 8082
CMD ["./main"]