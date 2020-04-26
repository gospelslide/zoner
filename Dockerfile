FROM golang:1.13

COPY . /app
WORKDIR /app

ARG 

RUN go build -o main .

CMD ["./main"]