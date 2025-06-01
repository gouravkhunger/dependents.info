# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git upx
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go && \
  upx --best --lzma main

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 5000
CMD ["./main"]
