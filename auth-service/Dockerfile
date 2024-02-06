# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Mateja Rilak"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./auth-service/go.mod ./auth-service/go.sum ./


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./auth-service/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
RUN mkdir data
COPY ./auth-service/data/blacklist-passwords.txt ./data/
COPY ./auth-service/security/rbac.json ./security/
RUN mkdir logs
COPY ./auth-service/logs/log.log ./logs/
COPY ./auth-service/server.crt /etc/ssl/certs/server.crt
COPY ./auth-service/server.csr /etc/ssl/csr/server.csr
COPY ./auth-service/server.key /etc/ssl/private/server.key
COPY ./auth-service/server.pem /etc/ssl/private/server.pem

# Command to run the executable
CMD ["./main"]