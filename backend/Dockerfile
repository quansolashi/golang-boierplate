# ビルド用コンテナ
FROM golang:1.23.0-alpine as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

WORKDIR /go/src/github.com/quansolashi/golang-boierplate/backend

RUN apk add --update --no-cache \
      git \
      make

COPY . .
RUN make build

# 実行用コンテナ
FROM alpine

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

RUN apk add --update --no-cache \
      ca-certificates \
      tzdata

COPY --from=builder /go/src/github.com/quansolashi/golang-boierplate/backend/app /app

EXPOSE 8080
CMD ["./app"]
