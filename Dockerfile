FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /http-helper
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hh cmd/hh/main.go

FROM alpine
WORKDIR /
RUN mkdir /data
COPY --from=builder /http-helper/hh /hh
COPY .env .env
EXPOSE 8000

CMD ["/hh", "-bind", "0.0.0.0", "-directory", "/data"]