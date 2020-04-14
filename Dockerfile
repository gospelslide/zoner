FROM golang:1.13

COPY . /app
WORKDIR /app

RUN go build -o main .

RUN ["./main"]