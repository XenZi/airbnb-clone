FROM golang:latest as builder

WORKDIR /app

COPY ./auth-service/go.mod ./auth-service/go.sum ./
COPY ./api-gateway ../api-gateway
RUN go mod download

COPY ./auth-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]