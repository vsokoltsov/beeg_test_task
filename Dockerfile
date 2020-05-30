FROM golang:1.13-alpine as binaries

RUN apk add git gcc g++

RUN go get github.com/cespare/reflex && \ 
    go get -u github.com/pressly/goose/cmd/goose

FROM golang:1.13-alpine

RUN apk add apache2-utils

COPY --from=binaries /go/bin /go/bin
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o app .