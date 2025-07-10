FROM golang:alpine AS Build
WORKDIR /socket

COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY backend ./
RUN go build -o app ./cmd

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /socket/app .
CMD ["app"]