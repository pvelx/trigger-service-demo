FROM golang:1.15

WORKDIR /app

COPY ./ /app

RUN go mod download

ENTRYPOINT go run main.go