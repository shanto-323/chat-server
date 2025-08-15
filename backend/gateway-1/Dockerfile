FROM golang:alpine AS build
WORKDIR /backend/gateway-1

COPY ./gateway-1/go.mod ./gateway-1/go.sum ./
RUN go mod download

COPY ./gateway-1 .

RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /backend/gateway-1/app .
CMD ["app"]