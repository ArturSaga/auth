FROM golang:1.23.2-alpine AS builder

COPY . /github.com/ArturSaga/auth/source/
WORKDIR /github.com/ArturSaga/auth/source/

RUN go mod download
RUN go build -o ./bin/auth-service cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ArturSaga/auth/source/bin/auth-service .

CMD ["./auth-service"]