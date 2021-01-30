FROM golang:1.15

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN make build

CMD ["/app/main"]