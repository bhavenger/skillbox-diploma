FROM golang:1.16.6-alpine3.14 as builder
RUN mkdir /build
WORKDIR /build
COPY . . /build/

RUN go mod download
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/app.go

FROM alpine:3.14
COPY --from=builder /build/app .

EXPOSE 8080

CMD ["./app"]