# syntax=docker/dockerfile:1
FROM golang:1.18-alpine as build
WORKDIR /app

RUN apk add build-base

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./bin/aiworkclient -buildvcs=false

FROM alpine:3
WORKDIR /app

ARG BUILD_ID="22.2.2"
ENV BUILD_ID=$BUILD_ID

COPY --from=build /app/configs.env ./secrets/configs.env
COPY --from=build /app/bin/aiworkclient ./aiworkclient

EXPOSE 8080
ENTRYPOINT ["/app/aiworkclient"]