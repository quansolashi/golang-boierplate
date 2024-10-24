FROM golang:1.22.0-alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

ENV AIR_VERSION 1.41.0

WORKDIR /go/src/github.com/EdvFuture/edvpath/backend

RUN apk add --update \
      git \
      make

RUN go install github.com/cosmtrek/air@v${AIR_VERSION}
