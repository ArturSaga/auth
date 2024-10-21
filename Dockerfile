FROM golang:1.23.2-alpine AS builder

COPY . /github.com/ArturSaga/auth/source/
WORKDIR /github.com/ArturSaga/auth/source/

RUN go mod download
RUN go build -o ./bin/auth-service-prod cmd/main.go --config-path=prod.env
RUN go build -o ./bin/auth-service-local cmd/main.go --config-path=local.env

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ArturSaga/auth/source/bin/auth-service-prod .
COPY --from=builder /github.com/ArturSaga/auth/source/bin/auth-service-local .

CMD ["./auth-service-prod"]
CMD ["./auth-service-local"]