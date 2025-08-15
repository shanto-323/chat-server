FROM golang:alpine AS build
WORKDIR /backend/client-service

COPY ./client-service/go.mod ./client-service/go.sum ./
RUN go mod download

COPY ./client-service . 

RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /backend/client-service/app .
CMD ["app"]