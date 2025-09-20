FROM golang:alpine AS build
WORKDIR /backend/message-service

COPY ./message-service/go.mod ./message-service/go.sum ./
RUN go mod download

COPY ./message-service . 

RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /backend/message-service/app .
CMD ["app"]
