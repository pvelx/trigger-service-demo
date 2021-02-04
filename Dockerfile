FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY . .
RUN go mod download
RUN go mod verify

RUN GOOS=linux GOARCH=amd64 go build -o /app/main github.com/pvelx/triggerServiceDemo

FROM alpine
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/main /app/main
USER appuser:appuser
ENTRYPOINT ["/app/main"]