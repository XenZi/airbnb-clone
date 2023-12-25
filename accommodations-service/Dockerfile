# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Aleksandar Pavlovic"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./accommodations-service/go.mod ./accommodations-service/go.sum ./


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./accommodations-service/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .


COPY ./accommodations-service/security/rbac.json ./security/

# Command to run the executable
CMD ["./main"]