FROM golang:1.17 AS build

ADD . /app
WORKDIR /app
RUN go mod tidy
RUN go build ./cmd/main.go