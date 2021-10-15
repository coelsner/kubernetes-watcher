# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /watcher

##
## Deploy
##
FROM alpine:latest

WORKDIR /

VOLUME /.kube/config

COPY --from=build /watcher /watcher

ENTRYPOINT ["/watcher"]