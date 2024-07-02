FROM golang:1.21-alpine AS builder

WORKDIR /automate/test
COPY . .

ENV GO111MODULE=on

RUN go mod tidy && go mod vendor
RUN cd app && go build -o main

FROM alpine:latest

RUN apk update && apk add curl

COPY --from=builder /automate/test /automate/test

WORKDIR /automate/test