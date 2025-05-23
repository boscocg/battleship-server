FROM golang:1.23.7 AS builder

WORKDIR /app

COPY . .

ARG ENV_FILE
COPY ${ENV_FILE} .env
COPY cloudbuild.yaml ./

COPY go.mod go.sum ./
RUN go mod download

# Build the application
RUN GOOS=linux GOARCH=amd64 go build -o battledak-server cmd/main.go

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
COPY --from=builder /app/${ENV_FILE} .env

EXPOSE 8080

CMD ["./battledak-server"]
