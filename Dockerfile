FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates curl wget
COPY --from=builder /app/server .
RUN adduser -D app && chown app:app /app/server
USER app
EXPOSE 5000
CMD ["./server"]