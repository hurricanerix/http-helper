FROM golang:1.24.2-bookworm AS build
WORKDIR /http-helper
COPY . .
RUN make bin/linux-amd64/hs

FROM alpine
WORKDIR /
RUN mkdir -p /data
COPY --from=build /http-helper/bin/linux-amd64/hs /bin/hs
COPY example.env .env
EXPOSE 8000

CMD ["/bin/hs", "-bind", "0.0.0.0", "-d", "/data"]