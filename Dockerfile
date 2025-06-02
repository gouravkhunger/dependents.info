# Step 1: Build Vite Frontend
FROM node:22-alpine AS vite-builder
WORKDIR /app

COPY www/package.json www/package-lock.json ./
RUN npm ci

COPY www/ ./
RUN npm run build

# Step 2: Build Go API
FROM golang:1.24-alpine AS go-builder
WORKDIR /app

RUN apk add --no-cache git upx

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ ./
COPY --from=vite-builder /app/dist ./static

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api main.go && \
  upx --best --lzma api

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=go-builder /app/api .
EXPOSE 5000
CMD ["./api"]
