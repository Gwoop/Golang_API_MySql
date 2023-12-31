# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /docker-gs-ping

EXPOSE 8000

CMD [ "/docker-gs-ping" ]