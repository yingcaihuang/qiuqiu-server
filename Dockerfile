FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o qiuqiu-server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/qiuqiu-server .

EXPOSE 8080

CMD ["./qiuqiu-server", "--addr", "0.0.0.0:8080", "--data", "/data"]
