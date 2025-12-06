# --- Builder stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary static để chạy trên Alpine
ENV CGO_ENABLED=0
RUN go build -o /go-shorty ./cmd/server/main.go

# --- Final image ---
FROM alpine:latest

# Cài CA certificates và tzdata để timezone nhận diện được
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

ENV PORT=8080
# Thiết lập timezone mặc định
ENV TZ=Asia/Ho_Chi_Minh

COPY --from=builder /go-shorty /usr/local/bin/go-shorty

EXPOSE 8080
CMD ["/usr/local/bin/go-shorty"]
