FROM golang:1.21.1

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]