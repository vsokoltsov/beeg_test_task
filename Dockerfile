FROM golang:1.13-alpine

RUN apk add git gcc g++ apache2-utils && \
    go get github.com/cespare/reflex && \ 
    go get -u github.com/pressly/goose/cmd/goose

WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o app .