FROM golang:1.22-alpine AS builder

ENV GOOS=darwin
ENV GOARCH=amd64

WORKDIR /go/src

COPY revshell.go .

RUN go build -o /go/bin/revshell-docker-${GOOS}-${GOARCH} revshell.go
