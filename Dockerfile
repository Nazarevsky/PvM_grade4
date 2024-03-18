FROM golang:latest

WORKDIR /app

COPY server ./

RUN go mod download

RUN go build main.go