# --- Builder stage ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /go-shorty ./cmd/server/main.go

# --- Final image ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /go-shorty ./go-shorty
RUN chmod +x ./go-shorty

RUN adduser -D -g '' appuser
USER appuser

ENV PORT=8080
ENV TZ=Asia/Ho_Chi_Minh

EXPOSE 8080
CMD ["./go-shorty"]