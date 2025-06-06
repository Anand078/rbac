FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o rbac-system cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rbac-system .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./rbac-system"]