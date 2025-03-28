FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ARG ENV_FILE
COPY . .
COPY ${ENV_FILE} ./

# Build the application
RUN GOOS=linux GOARCH=amd64 go build -o battledak-server cmd/api/main.go

FROM ubuntu:22.04

# Update and install system certificates
RUN apt-get update && \
    apt-get install -y ca-certificates curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/battledak-server .
RUN chmod +x /app/battledak-server

ARG ENV_FILE
COPY --from=builder /app/${ENV_FILE} ${ENV_FILE}

EXPOSE 8080

CMD ["./battledak-server"]
