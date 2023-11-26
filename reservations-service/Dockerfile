FROM golang:alpine as build_container

LABEL maintainer="Janko Rakonjac"

WORKDIR /app

COPY ./reservations-service/go.mod ./reservations-service/go.sum ./

RUN go mod download

COPY ./reservations-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build_container /app/main /


CMD ["./main"]