# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Boris Laus"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./user-service/go.mod ./user-service/go.sum ./


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./user-service/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
RUN mkdir logs
COPY ./user-service/logs/log.log ./logs/
COPY ./user-service/server.crt /etc/ssl/certs/server.crt
COPY ./user-service/server.csr /etc/ssl/csr/server.csr
COPY ./user-service/server.key /etc/ssl/private/server.key
COPY ./user-service/server.pem /etc/ssl/private/server.pem




# Command to run the executable
CMD ["./main"]