# FROM golang:1.26-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -o minjust-backend ./cmd

# FROM alpine:3.22

# WORKDIR /app

# COPY --from=builder /app/minjust-website .

# EXPOSE 8080

# CMD ["./minjust-website"]

FROM golang:1.26.3-alpine3.22

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o minjust-backend ./cmd


CMD ["./minjust-backend"]