FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go build -o weather ./cmd/app

FROM alpine:latest

RUN apk --no-cache add bash postgresql-client

WORKDIR /root/
COPY --from=builder /app/weather .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY ./migrations ./migrations
COPY ./entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]