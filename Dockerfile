FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o tax-calculator-cli ./main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/tax-calculator-cli .

ENTRYPOINT ["./tax-calculator-cli"]
