FROM golang:1.24.2-bookworm AS build
WORKDIR /http-helper
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make bin/linux-amd64/hs

FROM debian:stable-slim
WORKDIR /
RUN mkdir -p /data
COPY --from=build /http-helper/bin/linux-amd64/hs /bin/hs
COPY example.env .env
EXPOSE 8000

CMD ["/bin/hs", "-bind", "0.0.0.0", "-d", "/data"]