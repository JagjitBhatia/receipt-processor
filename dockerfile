# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY processor/ ./processor/

RUN go build -o /receipt-processor

EXPOSE 8080

CMD [ "/receipt-processor"]